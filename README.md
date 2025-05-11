
# Шина событий, работающая по принципу Publisher-Subscriber

## Требования к шине:

1. На один subject может подписываться (и отписываться) множество подписчиков.

2. Один медленный подписчик не должен тормозить остальных.

3. Нельзя терять порядок порядок сообщений (FIFO очередь).

4. Метод Close должен учитывать переданный контекст. Если он отменен - выходим сразу, работающие хендлеры оставляем работать.

5. Горутины (если они будут) течь не должны.

Ниже представлен API пакета subpub.

```sh
package subpub

import "context"

// MessageHandler is a callback function that processes messages delivered to subscribers.
type MessageHandler func(msg interface{})

type Subscription interface {
	// Unsubscribe will remove interest in the current subject subscription is for.
	Unsubscribe()
}

type SubPub interface {
	// Subscribe creates an asynchronous queue subscriber on the given subject.
	Subscribe(subject string, cb MessageHandler) (Subscription, error)

	// Publish publishes the msg argument to the given subject.
	Publish(subject string, msg interface{}) error

	// Close will shutdown sub-pub system.
	// May be blocked by data delivery until the context is canceled.
	Close(ctx context.Context) error
}

func NewSubPub() SubPub {
    panic("Implement me")
}
```

## Сервис подписок с использованием пакета subpub. Сервис работает по gRPC. Есть возможность подписаться на события по ключу и опубликовать события по ключу для всех подписчиков. 
## Protobuf-схема gRPC сервиса: 
```sh
import "google/protobuf/empty.proto";

syntax = "proto3";

service PubSub {
  // Подписка (сервер отправляет поток событий)
  rpc Subscribe(SubscribeRequest) returns (stream Event);

  // Публикация (классический запрос-ответ)
  rpc Publish(PublishRequest) returns (google.protobuf.Empty);
}

message SubscribeRequest {
  string key = 1;
}

message PublishRequest {
  string key = 1;
  string data = 2;
}

message Event {
  string data = 1;
}
```

## Как запускать сервис:

Клонируем репозиторий:
```sh
git clone https://github.com/V1Ro-Dev/pubSub
```

Переходим в папку config:
```sh
cd deploy/config
```

Заполняем файл server.yml, в котором необходимо указать:
```sh
addr:[Порт, на котором вы хотите запустить grpc сервер]
max_recv_msg_size: [Максимальный размер, получаемого сообщения]
max_send_msg_size: [Максимальный размер, отправляемого сообщения]
max_concurrent_streams: [Максимальное количество конкурентных потоков]
```

Переходим в корень проекта:
```sh
cd ../..
```

Запускаем приложение:
```sh
make up
```






