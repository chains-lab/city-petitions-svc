package problems

import (
	"context"
	"strconv"
	"time"

	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/meta"
	"github.com/chains-lab/city-petitions-svc/internal/constant"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InternalError(
	ctx context.Context,
) error {
	requestID := meta.RequestID(ctx)
	if requestID == uuid.Nil.String() {
		requestID = "unknown"
	}

	st := status.New(codes.Internal, "internal server error")

	info := &errdetails.ErrorInfo{
		Reason: "INTERNAL",
		Domain: constant.ServiceName,
		Metadata: map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	ri := &errdetails.RequestInfo{
		RequestId: requestID,
	}

	st, err := st.WithDetails(info, ri)
	if err != nil {
		return st.Err()
	}

	return st.Err()
}

func InvalidArgumentError(
	ctx context.Context,
	message string,
	violations ...*errdetails.BadRequest_FieldViolation,
) error {
	requestID := meta.RequestID(ctx)
	if requestID == uuid.Nil.String() {
		requestID = "unknown"
	}

	st := status.New(codes.InvalidArgument, message)

	info := &errdetails.ErrorInfo{
		Reason: canonicalString(st.Code()),
		Domain: constant.ServiceName,
		Metadata: map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	br := &errdetails.BadRequest{FieldViolations: violations}

	ri := &errdetails.RequestInfo{
		RequestId: requestID,
	}

	st, err := st.WithDetails(info, br, ri)
	if err != nil {
		return st.Err()
	}

	return st.Err()
}

func UnauthenticatedError(
	ctx context.Context,
	message string,
) error {
	requestID := meta.RequestID(ctx)
	if requestID == uuid.Nil.String() {
		requestID = "unknown"
	}

	st := status.New(codes.Unauthenticated, message)

	info := &errdetails.ErrorInfo{
		Reason: canonicalString(st.Code()),
		Domain: constant.ServiceName,
		Metadata: map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	ri := &errdetails.RequestInfo{
		RequestId: requestID,
	}

	st, err := st.WithDetails(info, ri)
	if err != nil {
		return st.Err()
	}

	return st.Err()
}

func PermissionDeniedError(
	ctx context.Context,
	message string,
) error {
	requestID := meta.RequestID(ctx)
	if requestID == uuid.Nil.String() {
		requestID = "unknown"
	}

	st := status.New(codes.PermissionDenied, message)

	info := &errdetails.ErrorInfo{
		Reason: canonicalString(st.Code()),
		Domain: constant.ServiceName,
		Metadata: map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	ri := &errdetails.RequestInfo{
		RequestId: requestID,
	}

	st, err := st.WithDetails(info, ri)
	if err != nil {
		return st.Err()
	}

	return st.Err()
}

func canonicalString(c codes.Code) string {
	switch c {
	case codes.OK:
		return "OK"
	case codes.Canceled:
		return "CANCELLED"
	case codes.Unknown:
		return "UNKNOWN"
	case codes.InvalidArgument:
		return "INVALID_ARGUMENT"
	case codes.DeadlineExceeded:
		return "DEADLINE_EXCEEDED"
	case codes.NotFound:
		return "NOT_FOUND"
	case codes.AlreadyExists:
		return "ALREADY_EXISTS"
	case codes.PermissionDenied:
		return "PERMISSION_DENIED"
	case codes.ResourceExhausted:
		return "RESOURCE_EXHAUSTED"
	case codes.FailedPrecondition:
		return "FAILED_PRECONDITION"
	case codes.Aborted:
		return "ABORTED"
	case codes.OutOfRange:
		return "OUT_OF_RANGE"
	case codes.Unimplemented:
		return "UNIMPLEMENTED"
	case codes.Internal:
		return "INTERNAL"
	case codes.Unavailable:
		return "UNAVAILABLE"
	case codes.DataLoss:
		return "DATA_LOSS"
	case codes.Unauthenticated:
		return "UNAUTHENTICATED"
	default:
		return "CODE(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}
