package labsession

import (
	"appseclabsplataform/config"
	"appseclabsplataform/database"
	labcluster "appseclabsplataform/services/labCluster"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"log/slog"

	"github.com/redis/go-redis/v9"
)

type LabSessionService struct {
	redisClient *redis.Client
	database    *database.Database
	labCluster  *labcluster.LabClusterService
	ctx         context.Context
}

func NewLabSessionService(database *database.Database, labCluster *labcluster.LabClusterService, config *config.Config) *LabSessionService {
	redisDB, err := strconv.Atoi(config.LabSessionServiceConfig.DB)
	if err != nil {
		panic(fmt.Sprintf("Invalid Redis DB value: %v", err))
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.LabSessionServiceConfig.Addr,
		Password: config.LabSessionServiceConfig.Password,
		DB:       redisDB,
	})

	return &LabSessionService{
		redisClient: redisClient,
		database:    database,
		labCluster:  labCluster,
		ctx:         context.Background(),
	}
}

func (s *LabSessionService) SetLabSession(session *LabSession, expiration time.Duration) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = s.redisClient.Set(s.ctx, fmt.Sprintf("lab_session:%s", session.Namespace), data, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *LabSessionService) GetLabSession(namespace string) (*LabSession, error) {
	data, err := s.redisClient.Get(s.ctx, fmt.Sprintf("lab_session:%s", namespace)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("lab session not found for namespace: %s", namespace)
		}
		return nil, err
	}

	var session LabSession
	err = json.Unmarshal([]byte(data), &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
func (s *LabSessionService) DeleteLabSession(namespace string) error {
	err := s.redisClient.Del(s.ctx, fmt.Sprintf("lab_session:%s", namespace)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *LabSessionService) ListenForLabSessionExpiry() {
	pubsub := s.redisClient.Subscribe(s.ctx, "__keyevent@0__:expired")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		if msg.Payload == "" {
			continue
		}

		var namespace string
		fmt.Sscanf(msg.Payload, "lab_session:%s", &namespace)

		err := s.labCluster.DeleteLabSession(namespace)
		if err != nil {
			slog.Error("Error deleting lab session from lab cluster %s: %v", namespace, err)
			continue
		}

		err = s.database.SetLabAttemptAsTimeout(namespace)
		if err != nil {
			slog.Error("Error setting lab attempt as timeout %s: %v", namespace, err)
		}

		err = s.DeleteLabSession(namespace)
		if err != nil {
			slog.Error("Error deleting lab session %s: %v", namespace, err)
		}
	}
}
