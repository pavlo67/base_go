# Правила оформлення інтерфейсів та імплементацій

Приклад пакету з інтерфейсом та імплементаціями: github.com/pavlo67/base_go/entities/files

# Опис типів даних та інтерфейсу

Опис типів даних та інтерфейсів пишемо в файлі: entities/files/operator.go

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

Імплементацію зберігаємо в каталозі: entities/files/files_<impl>

Приклад: entities/files/files_sqlite

В назвах файлів імплементації використовуємо назву каталогу як префікс (для ясности в назвах табів в IDE): entities/files/filesimpl_/filesimpl_....go

Використовуємо наступні службові пакети:

	"github.com/pavlo67/base_go/lib/db"     // інтерфейс для роботи з тестовою БД
	"github.com/pavlo67/base_go/lib/errors" // імплементація помилок
	"github.com/pavlo67/base_go/lib/logger" // логер

При ініціяції "структури-імплементації" ініціюємо логер-змінну l. 

При визначенні структури, яка імплементує певний інтерфейс, перевіряємо, чи імплементація коректна (якщо ні, цей рядок не скомпілюється). Наприклад:

    var _ files.Operator = &filesSQLite{}

    type filesSQLite struct {
        db                                             *sql.DB
        stmSave, stmRead, stmRemove, stmList, stmClean *sql.Stmt
    }

Кожну помилку, яка приходить з викликів до службових функцій (зокрема, це стосується викликів стандартних бібліотек), 
обгортаємо за допомогою константи, яка має бути визначена для кожного нашого методу. Локально створені помилки також формуємо з використанням цієї константи, 
щоб не дублювати назву методу в текстах помилок.

Всі параметри в функцію ініціяції передаються явно, не через змінні оточення. Приклад функції ініціяції:
    
    var l logger.Operator

    const onNew = "on files_sqlite.Init():"
    
    func New(dsn string, create bool, l_ logger.Operator) (files.Operator, db.Operator, error) {
        if l_ == nil {
    		return nil, nil, errors.New("", onNew + " l_ == nil")
        }
        l = l_
    
        sqlDB, err := sql.Open("sqlite3", dsn)
        if err != nil {
            return nil, nil, errors.Wrap(err, onNew)
        }
    
        op := &filesSQLite{db: sqlDB}
    
        if create {
            if err = op.Create(sqlDB); err != nil {
                return nil, nil, errors.Wrap(err, onNew)
            }
        }
    
        sqlStmts := []sqllib.SqlStmt{
            {&op.stmSave, sqlSave},
            {&op.stmRead, sqlRead},
            {&op.stmRemove, sqlRemove},
            {&op.stmList, sqlList},
            {&op.stmClean, sqlClean},
        }
    
        for _, sqlStmt := range sqlStmts {
            if err := sqllib.PrepareQuery(sqlDB, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
                return nil, nil, errors.Wrap(err, onNew)
            }
        }
    
        return op, op, nil
	}


## Особливості SQL-імплементацій 

Використовуємо наступний службовий пакет:

	"github.com/pavlo67/base_go/lib/sqllib"     // інтерфейс для роботи з SQL БД

Адресу БД не зберігаємо в імплементації, вона повинна братись з конфігів (назовні цієї імплементації) і передаватись як параметр в функцію New().

Таблиці, необхідні для роботи, створюємо за умови, що параметр create == true (це використовується виключно при тестах). SQL-запит, необхідний для створенян таблиць,
тримаємо в файлі create.sql.

Кожен SQL-запит тримаємо у відповідній константі.

Для підготовки БД в функції New() береться  *sql.DB і набір SQL-запитів і робиться набір "препарованих запитів", які потім будуть, за потреби, виконуватись. Основна ідея — 
відчасти, швидкодія імплементації але, головне, в момент ініціяції інтерфейсу таким чином перевіряється валідність бази даних. Якщо нема таблиці, або полів в таблиці, або ще 
щось не так — система звалиться на старті з явною діягностикою і це набагато зручніше, ніж ловити потім runtime-помилки. Приклад — вище.

Базовий DB-інтерфейс містить методи, які використовуються тільки в тестовому оточенні (перевірка оточення повинна робитись в кожній імплементації кожного з цих методів):

    Create(db *sql.DB) error
    Clean() error


# Тести

Тестовий сценарій (це не юніт-тести!) зберігається в файлі entities/files/operator_test_scenario.go і має бути незалежним від імплементацій. 

На прикладі files його можна, в загальних словах, описати наступним чином: 

    import "github.com/stretchr/testify/require" 

    func FilesTestScenario(t *testing.T, filesOp files.Operator, filesDB db.Operator) { 
        err := filesDB.Clean() 
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

Повний код сценарію: entities/files/operator_test_scenario.go

Пускач тестів створюється в кожній імплементації — він повинен створити логер, ініціювати основний інтерфейс і викликати тестовий сценарій.
        
Приклад пускача: entities/files/files_sqlite/files_sqlite_test.go 



