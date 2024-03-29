package interceptors

import (
	"context"
	"github.com/Mldlr/storety/internal/client/config"
	"github.com/Mldlr/storety/internal/client/pkg/utils"
	"github.com/Mldlr/storety/internal/constants"
	pb "github.com/Mldlr/storety/internal/proto"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

// RetryClientInterceptor is a client interceptor that retries RPCs.
type RetryClientInterceptor struct {
	client        pb.UserClient
	cfg           *config.Config
	retryTimes    uint
	retryDuration time.Duration
}

// NewRetryClientInterceptor creates a new RetryClientInterceptor and returns a pointer to it.
// It takes a configuration object, retryTimes, and retryDuration as parameters.
func NewRetryClientInterceptor(cfg *config.Config, retryTimes uint, retryDuration time.Duration, conn *grpc.ClientConn) *RetryClientInterceptor {
	return &RetryClientInterceptor{
		client:        pb.NewUserClient(conn),
		cfg:           cfg,
		retryTimes:    retryTimes,
		retryDuration: retryDuration,
	}
}

// UnaryInterceptor is a modified grpc-middleware unary interceptor that retries RPCs and refreshes the token if needed.
func (r *RetryClientInterceptor) UnaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, callOpts ...grpc.CallOption) error {
	var lastErr error
	for attempt := uint(0); attempt < r.retryTimes; attempt++ {
		if err := waitRetryBackoff(ctx, attempt); err != nil {
			return err
		}
		lastErr = invoker(ctx, method, req, reply, cc)
		if lastErr == nil {
			return nil
		}
		if isContextError(lastErr) {
			if ctx.Err() != nil {
				return lastErr
			}
			continue
		}
		if rpctypes.ErrorDesc(lastErr) == constants.ErrExpiredToken.Error() {
			log.Println("Token expired, trying to refresh token")
			request := &pb.RefreshUserSessionRequest{}
			result, err := r.client.RefreshUserSession(ctx, request)
			if err != nil {
				return err
			}
			r.cfg.UpdateTokens(result.AuthToken, result.RefreshToken)
			continue
		}
		break
	}
	return lastErr
}

func waitRetryBackoff(ctx context.Context, attempt uint) error {
	waitTime := time.Duration(0)
	if attempt > 0 {
		waitTime = utils.JitterUp(50*time.Millisecond /*jitter*/, 0.10)
	}
	if waitTime > 0 {
		timer := time.NewTimer(waitTime)
		select {
		case <-ctx.Done():
			timer.Stop()
			return contextErrToGrpcErr(ctx.Err())
		case <-timer.C:
		}
	}
	return nil
}

func isContextError(err error) bool {
	return status.Code(err) == codes.DeadlineExceeded || status.Code(err) == codes.Canceled
}

func contextErrToGrpcErr(err error) error {
	switch err {
	case context.DeadlineExceeded:
		return status.Errorf(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return status.Errorf(codes.Canceled, err.Error())
	default:
		return status.Errorf(codes.Unknown, err.Error())
	}
}
