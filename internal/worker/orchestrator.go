package worker

type ArcDownTask struct {
	UUID       string
	TaskID     int64
	Expected   int64
	Done       int64
	Collection string
	IsArchived bool
}

type TaskDone struct {
	UUID   string
	TaskID int64
}

const (
	TaskQuit int64 = iota
	TaskDownload
	TaskArchive
)

func (man *WorkerManager) DownloadAndArchive(uuid string, collection string, sum int64) {
	man.TaskChannel <- ArcDownTask{
		UUID:       uuid,
		TaskID:     TaskDownload,
		Expected:   sum,
		Done:       0,
		Collection: collection,
		IsArchived: true,
	}
}

// Orchaestrator function to handle archiver and downloader
func (man *WorkerManager) Orchestrator() {
	ongoingTask := map[string]*ArcDownTask{}

	for {
		select {
		case o := <-man.TaskChannel:
			if o.TaskID == TaskQuit {
				if len(ongoingTask) != 0 {
					man.TaskChannel <- o
					break
				} else {
					close(man.TaskChannel)
					close(man.DoneChannel)
					return
				}
			}

			if o.TaskID == TaskDownload {
				ongoingTask[o.UUID] = &ArcDownTask{
					UUID:       o.UUID,
					TaskID:     o.TaskID,
					Expected:   o.Expected,
					Done:       0,
					IsArchived: true,
					Collection: o.Collection,
				}
				break
			}

			if o.TaskID == TaskArchive {
				ongoingTask[o.UUID] = &ArcDownTask{
					UUID:     o.UUID,
					TaskID:   o.TaskID,
					Expected: 1,
					Done:     0,
				}
			}

		case done := <-man.DoneChannel:
			if meta, ok := ongoingTask[done.UUID]; ok {
				meta.Done += 1
			}
		}

		for k, task := range ongoingTask {
			if task.Done == task.Expected {
				delete(ongoingTask, k)
			}
		}
	}
}
