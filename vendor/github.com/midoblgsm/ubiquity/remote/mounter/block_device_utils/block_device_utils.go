package block_device_utils

import (
	"github.com/midoblgsm/ubiquity/utils"
	"github.com/midoblgsm/ubiquity/utils/logs"
)

type blockDeviceUtils struct {
	logger logs.Logger
	exec   utils.Executor
}

func NewBlockDeviceUtils() BlockDeviceUtils {
	return newBlockDeviceUtils(utils.NewExecutor())
}

func NewBlockDeviceUtilsWithExecutor(executor utils.Executor) BlockDeviceUtils {
	return newBlockDeviceUtils(executor)
}

func newBlockDeviceUtils(executor utils.Executor) BlockDeviceUtils {
	return &blockDeviceUtils{logger: logs.GetLogger(), exec: executor}
}
