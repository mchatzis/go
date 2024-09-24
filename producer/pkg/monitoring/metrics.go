package monitoring

import (
	"context"
	"strconv"
	"time"

	"github.com/mchatzis/go/producer/pkg/logging"
	"github.com/mchatzis/go/producer/pkg/sqlc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var logger = logging.GetLogger()

var taskTypes = [10]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

var totalValueOfDoneTasksByType *prometheus.GaugeVec
var countOfDoneTasksByType *prometheus.GaugeVec
var countOfTasksByState *prometheus.GaugeVec

var queries *sqlc.Queries
var taskStates = []sqlc.TaskState{
	sqlc.TaskStatePending,
	sqlc.TaskStateProcessing,
	sqlc.TaskStateDone,
	sqlc.TaskStateFailed,
}

func RegisterCollectors() {
	promauto.NewGauge(prometheus.GaugeOpts{
		Name: "service_up",
		Help: "Indicates whether the service is up (1) or down (0)",
	}).Set(1)

	totalValueOfDoneTasksByType = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "total_value_done_tasks_by_type",
		Help: "Sums all the values of done tasks for each type",
	}, []string{"type"})

	countOfDoneTasksByType = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "count_done_tasks_by_type",
		Help: "Counts how many tasks have been processed for each type",
	}, []string{"type"})

	for _, t := range taskTypes {
		totalValueOfDoneTasksByType.With(prometheus.Labels{"type": t})
		countOfDoneTasksByType.With(prometheus.Labels{"type": t})
	}

	countOfTasksByState = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "count_tasks_by_state",
		Help: "Counts how many tasks exist for each state",
	}, []string{"state"})

	for _, state := range taskStates {
		countOfTasksByState.With(prometheus.Labels{"state": string(state)})
	}

}

func UpdateMetrics(queriesIn *sqlc.Queries) {
	queries = queriesIn
	for {
		updateTotalValueOfDoneTasksByTypeGauges()
		updateCountOfDoneTasksByTypeGauge()
		updateCountOfTasksByStateGauge()
		time.Sleep(100 * time.Millisecond)
	}
}

func updateTotalValueOfDoneTasksByTypeGauges() {
	rows, err := queries.GetTotalValueOfDoneTasksByType(context.Background())
	if err != nil {
		logger.Errorf("Error in monitoring, updating total value of done tasks by type failed: %v", err)
	}

	for _, r := range rows {
		taskType := strconv.Itoa(int(r.Type))
		valueSum := r.Sum
		totalValueOfDoneTasksByType.With(prometheus.Labels{"type": taskType}).Set(float64(valueSum))
	}
}

func updateCountOfDoneTasksByTypeGauge() {
	rows, err := queries.GetCountOfDoneTasksByType(context.Background())
	if err != nil {
		logger.Fatalf("Error in monitoring, counting done tasks by type failed: %v", err)
	}

	for _, r := range rows {
		taskType := strconv.Itoa(int(r.Type))
		taskCount := r.Count
		countOfDoneTasksByType.With(prometheus.Labels{"type": taskType}).Set(float64(taskCount))
	}
}

func updateCountOfTasksByStateGauge() {
	rows, err := queries.GetCountOfTasksByState(context.Background())
	if err != nil {
		logger.Fatalf("Error in monitoring, counting tasks by state failed: %v", err)
	}

	for _, state := range taskStates {
		countOfTasksByState.With(prometheus.Labels{"state": string(state)}).Set(0)
		for _, r := range rows {
			if taskState := r.State; taskState == state {
				taskCount := r.Count
				countOfTasksByState.With(prometheus.Labels{"state": string(taskState)}).Set(float64(taskCount))
			}
		}
	}

}
