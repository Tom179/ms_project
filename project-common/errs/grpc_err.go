package errs //用于转换错误信息
import (
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	common "test.com/project-common"
)

func GrpcError(err *BError) error {
	return status.Error(codes.Code(err.Code), err.Msg) //int转换未uint32
}

func ParseGrpcError(err error) (common.ResponseCode, string) {
	fromError, _ := status.FromError(err)
	return common.ResponseCode(fromError.Code()), fromError.Message()
}
