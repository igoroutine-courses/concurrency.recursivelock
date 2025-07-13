# Recursive lock

## Описание
На занятиях мы проговаривали недостатки мьютекса в Go

*Разблокировать мьютекс может горутина, которая им не владеет*
```go
var l sync.Mutex
var v string

func f() {
	v = "the nature of concurrency"
	l.Unlock()
}

func main() {
	l.Lock()
	go f()
	l.Lock()
	fmt.Println(v)
}
```

*Одна горутина не может взять Lock() два раза, это может доставлять боль при написании рекурсивных коллекций*
```go
func main() {
	r := Reentrant{
		mx: new(sync.Mutex),
	}

	r.Outer()
}

type Reentrant struct {
	mx *sync.Mutex
}

func (r *Reentrant) Outer() {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.Inner()
}

func (r *Reentrant) Inner() {
	r.mx.Lock()
	defer r.mx.Unlock()
}
```

Например, в *Java* у `ReentrantLock` таких недостатков нет. Такая реализация мьютекса традиционно называется
рекурсивной/реентерабельной (reentrant).

```sql
import java.util.concurrent.locks.ReentrantLock;

public class Main {
    public static void main(String[] args) {
        Reentrant r = new Reentrant();
        r.outer();

        System.out.println("All done!");
    }
}

class Reentrant {
    private final ReentrantLock lock = new ReentrantLock();

    public void outer() {
        lock.lock();

        try {
            inner();
        } finally {
            lock.unlock();
        }
    }

    public void inner() {
        lock.lock();

        try {
            // action
        } finally {
            lock.unlock();
        }
    }
}
```

В этом задании вам необходимо написать подобную реализацию на Go

```go
// ErrUnlockFromAnotherGoroutine - ошибка, если Unlock вызывает горутина, не владевшая мьютексом
var ErrUnlockFromAnotherGoroutine = errors.New("unlock from non-owner goroutine")

// New создаёт новый экземпляр рекурсивного (реентерабельного) мьютекса
func New() *Mutex {
	panic("not implemented")
}

// Mutex реализует реентрантный мьютекс
type Mutex struct{}

// Lock захватывает мьютекс. Если текущая горутина уже владеет им,
// повторный вызов безопасен (реентерабельность)
func (r *Mutex) Lock() {
	panic("not implemented")
}

// Unlock освобождает мьютекс. Только горутина-владелец может вызвать Unlock, если другая горутина
// попытается разблокировать - будет паника
func (r *Mutex) Unlock() {
	panic("not implemented")
}
```

Больше всего такой мьютекс пригождается в рекурсивных реализациях, например, в персистентном дереве конфигураций

```go
// ErrEffectiveValueNotFound
// Ошибка, возвращаемая, если ни в текущем узле, ни у предков не найдено значение
var ErrEffectiveValueNotFound = errors.New("effective value not found")

// Node представляет узел конфигурационного дерева.
// Использует обобщения (type parameter T), поддерживает потокобезопасность
type Node[T any] struct{}

// NewNode создает новый узел конфигурационного дерева.
// Если передан родитель, добавляет новый узел в список его детей
func NewNode[T any](parent *Node[T], value T) *Node[T] {
	panic("not implemented")
}

// ClearValue очищает локальное значение в текущем узле
func (n *Node[T]) ClearValue() {
	panic("not implemented")
}

// All возвращает все предыдущие узлы
func (n *Node[T]) All() iter.Seq[*Node[T]] {
	panic("not implemented")
}

// GetEffectiveValue возвращает эффективное значение текущего узла.
// Если локальное значение отсутствует, рекурсивно ищет значение у предков
func (n *Node[T]) GetEffectiveValue() (T, error) {
	panic("not implemented")
}
```

## Задание

### Часть 1. Написание мьютекса
Напишите реализацию  мьютекса в [lock.go](./internal/reentrant/lock.go) с помощью библиотеки [goid](https://github.com/petermattis/goid). Вы можете проверить корректность
вашего решения, используя тесты [lock_test.go](./internal/reentrant/lock_test.go), которые эмулируют проблемы стандартного мьютекса Go.

### Часть 2. Использование мьютекса
Потренеруйтесь использовать мьютекс в Go. Для этого напишите потокобезопасную реализацию персистентного дерева конфигураций в
[config.go](./internal/config/config.go) и проверьте её корретность, запустив тесты
[config_test.go](./internal/config/config_test.go). Об этом дереве можно думать, как о персистентном хранилище, в котором
можно получать актуальные значения `GetEffectiveValue` и удалять неактуальные `ClearValue`, а также получать историю значений `All`.

## Сдача
* Создать ветку с решением, открыть pull request из этой ветки в ветку `main` **вашего репозитория**.

* Если у вас есть ревью от преподавателя, отправить [в задании на GetCourse](https://igoroutine.getcourse.ru/pl/teach/control/lesson/view?id=342566686&editMode=0) ссылку на PR.

## Запуск проверки
Для запуска скриптов на курсе необходимо установить [go-task](https://taskfile.dev/installation/)

```
go install github.com/go-task/task/v3/cmd/task@latest
```

Запустить тесты:

```bash 
task test
```

Подтянуть изменения в репозиторий:

```bash
task update
```
