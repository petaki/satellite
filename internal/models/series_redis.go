package models

import "github.com/gomodule/redigo/redis"

type RedisSeriesRepository struct {
	RedisPool      *redis.Pool
	RedisKeyPrefix string
}

func (rsr *RedisSeriesRepository) FindCpu(seriesType SeriesType) (Series, error) {
	return nil, nil
}

func (rsr *RedisSeriesRepository) FindMemory(seriesType SeriesType) (Series, error) {
	return nil, nil
}

func (rsr *RedisSeriesRepository) FindDisk(seriesType SeriesType, path string) (Series, error) {
	return nil, nil
}

func (rsr *RedisSeriesRepository) FindDiskPaths(cursor string, paths []string) ([]string, error) {
	return nil, nil
}
