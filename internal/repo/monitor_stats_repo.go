package repo

import (
	"context"

	"github.com/dushixiang/pika/internal/models"
	"github.com/go-orz/orz"
	"gorm.io/gorm"
)

type MonitorStatsRepo struct {
	orz.Repository[models.MonitorStats, uint]
	db *gorm.DB
}

func NewMonitorStatsRepo(db *gorm.DB) *MonitorStatsRepo {
	return &MonitorStatsRepo{
		Repository: orz.NewRepository[models.MonitorStats, uint](db),
		db:         db,
	}
}

// FindByAgentAndName 根据探针ID和监控名称查找统计数据
func (r *MonitorStatsRepo) FindByAgentAndName(ctx context.Context, agentID, monitorName string) (*models.MonitorStats, error) {
	var stats models.MonitorStats
	err := r.db.WithContext(ctx).
		Where("agent_id = ? AND monitor_name = ?", agentID, monitorName).
		First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// UpsertStats 插入或更新统计数据
func (r *MonitorStatsRepo) UpsertStats(ctx context.Context, stats *models.MonitorStats) error {
	return r.db.WithContext(ctx).
		Where("agent_id = ? AND monitor_name = ?", stats.AgentID, stats.MonitorName).
		Assign(stats).
		FirstOrCreate(stats).Error
}

// ListByMonitorName 根据监控名称列出所有探针的统计数据
func (r *MonitorStatsRepo) ListByMonitorName(ctx context.Context, monitorName string) ([]models.MonitorStats, error) {
	var statsList []models.MonitorStats
	err := r.db.WithContext(ctx).
		Where("monitor_name = ?", monitorName).
		Find(&statsList).Error
	return statsList, err
}

// ListAll 列出所有统计数据
func (r *MonitorStatsRepo) ListAll(ctx context.Context) ([]models.MonitorStats, error) {
	var statsList []models.MonitorStats
	err := r.db.WithContext(ctx).Find(&statsList).Error
	return statsList, err
}
