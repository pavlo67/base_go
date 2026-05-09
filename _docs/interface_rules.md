# Правила оформлення інтерфейсів та імплементацій

Приклад пакету з інтерфейсом та імплементаціями: /<path_to_example>/<some>

# Опис типів даних та інтерфейсу

Приклад: https://github.com/pavlo67/base_go/tree/master/entities/files

Опис типів даних та інтерфейсів пишемо в файлі: <path_to_example>/<some>/operator.go

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


# Імплементація

Приклад: https://github.com/pavlo67/base_go/tree/master/entities/files/files_sqlite

Імплементацію зберігаємо в каталозі: <path_to_example>/<some>/<some>_<impl>/

В назвах файлів імплементації використовуємо назву каталогу як префікс (для ясности в назвах табів в IDE): <path_to_example>/some/some_<impl>/some_<impl>....go

Використовуємо наступні службові пакети:

	"github.com/pavlo67/base_go/lib/db"     // інтерфейс для роботи з тестовою БД
	"github.com/pavlo67/base_go/lib/errors" // імплементація помилок
	"github.com/pavlo67/base_go/lib/logger" // логер

При ініціяції "структури-імплементації" ініціюємо логер-змінну l. 

При визначенні структури, яка імплементує наш інтерфейс, перевіряємо, чи імплементація коректна (якщо ні, цей рядок не скомпілюється). Наприклад:

    var _ files.Operator = &filesSQLite{}

Кожну помилку, яка приходить з викликів до службових функцій (особливо це стосується до викликів стандартних бібліотек), 
обгортаємо за допомогою константи, яка визначається для кожного нашого методу. 

Сигнатура ініціяції (всі параметри передаються явно, не через змінні оточення): 

    func New(<params>, l_ logger.Operator) (<some>.Operator, db.Cleaner, error) {

Приклад сигнатури ініціяції (функція New()):

    type filesSQLite struct {
        db *sql.DB
    }
    
    var l logger.Operator

    const onNew = "on files_sqlite.Init():"
    
    func New(dsn string, l_ logger.Operator) (files.Operator, db.Cleaner, error) {
        if l_ == nil {
            return nil, nil, errors.New("", "l_ == nil")
        }
        l = l_

	    sqlDB, err := sql.Open("sqlite3", dsn)
	    if err != nil {
		    return nil, nil, errors.Wrap(err, onNew)
	    }

	    op := &filesSQLite{db: sqlDB}

        // db preparation is done here       

    	return op, op, nil
	}


## Особливості SQL-імплементацій 

Для підготовки БД в функції New() береться  *sql.DB і набір SQL-запитів і робиться набір "препарованих запитів", які потім будуть, за потреби, виконуватись. Основна ідея — 
відчасти, швидкодія імплементації але, головне, в момент ініціяції інтерфейсу таким чином перевіряється валідність бази даних. Якщо нема таблиці, або полів в таблиці, або ще 
щось не так — система звалиться на старті з явною діягностикою і це набагато зручніше, ніж ловити потім runtime-помилки.

# Тести

Тестовий сценарій (це не юніт-тести!) зберігається в файлі <path_to_example>/<some>/operator_test_scenario.go і описується наступним чином (на прикладі files): 

    import "github.com/stretchr/testify/require" 

    func FilesTestScenario(t *testing.T, filesOp files.Operator, filesDB db.Operator) { 
        err := filesCleaner.Clean() 
        require.NoError(t, err)} 

        filesOp.List(...) 
        require(<no items>) 

        // --------------------------------------------

        data1 := files.Data{...} 
        data2 := files.Data{...}
        data3 := files.Data{...}

        ... 

        filesOp.Save(data1) 
        filesOp.List(...) 
        require( <only item1 in list> ) 

        item1Readed, err1 := filesOp.Read(item1.Path) 
        require.NoError(err1) 
        require.NotNil(item1Readed) 
        require.Equal(data1, item1Readed.Data, <error explaining>) 

        filesOp.Save(data1)        // re-save 
        filesOp.List(...) 
        require( <only item1 in list> ) 

        // and more in the same style: 
        //      save data2; check list; read and check item2; read and check item1 (must be unchanged); 
        //      save data3; check list; read and check item3; read and check item2 (must be unchanged);.... 
        //      remove item2 and check list and check item1 and item3 are unchanged 
        // всі data в одному каталозі (з точки зору filepath.Dir),

        
// виклик .Clean() нехай буде захищений перевіркою ENV-змінної — чи це TEST-оточення, тільки в ньому можна чистити базу 
// test-файл, який ініціює filesSQLite і викликає FilesTestScenario()



