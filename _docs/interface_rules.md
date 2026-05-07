## Правила оформлення інтерфейсів та імплементацій

Приклад пакету з інтерфейсом та імплементаціями: /<path_to_example>/some

Опис інтерфейсу: /<path_to_example>/some/operator.go

Якщо є основна сутність, з якою працює інтерфейс (така сутність завжди наявна в CRUD-інтерфейсах), то вона зветься Data (назва однакова для всіх інтерфейсів):

    type Data struct {
        <поля сутності> 
    }

Якщо є поля, які додаються до основної сутности автоматично, то вони не включаються в Data struct, але формується структура-обгортка Item (назва однакова для всіх інтерфейсів). 
Зокрема, такими полями є CreatedAt та UpdatedAt (час створення і модифікації запису в БД):

    type Item struct {
        Data      `          json:",inline" bson:",inline"`
        CreatedAt time.Time `json:",omitzero" bson:",omitempty"`  
        UpdatedAt time.Time `json:",omitzero" bson:",omitempty"`
    }

Базовий CRUD-інтерфейс при цьому описується наступним чином:
    
    type Operator interface {
        // creates new or replaces existing Item's record
        Save(data Data) error
        Read(<unique key/id>) (*Item, error)
        Remove(<unique key/id>) error
        List(<conditions>) ([]Item, error)
    }

// перевірка типу інтерфейсу
// обгортання помилок


