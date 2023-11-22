package mocks

import (
	"config_center/api/types"
	"github.com/stretchr/testify/mock"
)

type MockDependenciesRepository struct {
	mock.Mock
}

func (r *MockDependenciesRepository) FindByMajor(table, platform string, major int) (types.Dependency, error) {
	args := r.Called(table, platform, major)
	return args.Get(0).(types.Dependency), args.Error(1)
}

func (r *MockDependenciesRepository) FindByMajorMinor(table, platform string, major, minor int) (types.Dependency, error) {
	args := r.Called(table, platform, major, minor)
	return args.Get(0).(types.Dependency), args.Error(1)
}

func (r *MockDependenciesRepository) FindByMajorMinorPatch(table, platform string, major, minor, patch int) (types.Dependency, error) {
	args := r.Called(table, platform, major, minor, patch)
	return args.Get(0).(types.Dependency), args.Error(1)
}