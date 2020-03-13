package middlewares

import (
	"building-microservices-with-go.com/logging/entities"
	"building-microservices-with-go.com/logging/httputil"
	"context"
	"encoding/json"
	"github.com/alexcesaro/statsd"
	"github.com/sirupsen/logrus"
	"net/http"
)

// ValidationMiddleware는 핸들러를 실행시키기 전에, 요청의 유효성을 검사하여 그 결과에 따라 처리를 하는 미들웨어이다.
type ValidationMiddleware struct {
	statsD  *statsd.Client
	logger  *logrus.Logger
	next 	http.Handler
}

func NewValidationMiddleware(statsD *statsd.Client, logger *logrus.Logger, next http.Handler) *ValidationMiddleware {
	return &ValidationMiddleware{statsD: statsD, logger: logger, next: next}
}

func (vm *ValidationMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request entities.HelloWorldRequest

	// 만약 디코딩에 실패하였다면(Request Body의 유효성 인정X) 아래 구문을 실행시킨다.
	err := decoder.Decode(&request)
	if err != nil {
		// statsd.Client.Increment 메서드를 이용하여 매개변수로 주어진 버킷의 카운트를 1씩 증가시킨다.
		vm.statsD.Increment(validationFailed)
		message := httputil.RequestSerializer{Request: r}
		vm.logger.WithFields(logrus.Fields{
			"group": "middleware",
			"segment": "validation",
			"outcome": http.StatusBadRequest,
		}).Info(message.ToJSON())

		http.Error(rw, "Bad Request", http.StatusBadRequest)
	}

	// 다른 핸들러 또는 미들웨어에서도 Request Body에 접근할 수 있도록 하기 위해 컨텍스트에 저장시켜논다.
	c := context.WithValue(r.Context(), "name", request.Name)
	r = r.WithContext(c)

	// 이번에는 유효성 검사가 성공하였으니 statsd.Client.Increment 메서드의 매개변수에 성공했다는 버킷을 전달한다.
	vm.statsD.Increment(validationSuccess)
	// 마지막으로 다음 http.Handle.ServeHTTP를 실행시키고 ValidationMiddleware를 마무리 시킨다.
	vm.next.ServeHTTP(rw, r)
}