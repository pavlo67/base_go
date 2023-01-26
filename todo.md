## crud

* Backup() method

## transport

* transpOp.Track()
* "контракти" транспорту:
    - приймач перевіряє, чи параметри відповідають типу задачі
    - відправник дає інфу приймачу, які типи задач очікуються у відповідь (і той це перевіряє)
* "внутрішня відправка", якщо домени відправника і отримувача співпадають
* "внутрішня відправка відповіді", якщо отримувач повертає її у вихідному запиті
* завжди модифікувати packs.Pack.Key:
    - або унікальним префіксом/постфіксом
    - або робити унікальний, а ориґінал переміщати в History
* описати-систематизувати опції і сценарії
* онлайн-роутинг (DNS!)
* при конектах (а навіть при старті?) — перевіряти "очікувані" формати пакетів
* sender.Operator.Handle, sender.Operator.History
* мережа довірених адресатів

### json_any

* розбивати вивід на окремі файли

## data, packs, tasks

* use default key's value searching data with empty .Key
* check viewer_key



## actor

* check if (sub-)prefix is unique 


## files_fs

* options: follow_symlinks / allow_access_beyond_root 


## records

* RBAC checks
* "form" URL - get from serverConfig
* нормально обробляти errs в хендлах
* .History test

## storage

* pagination


## tagger

* tests for count...


## flow

* pagination (min, max, amount, limited_by)
* приховування/показ прочитаного раніше
* flow_cleaner --> dataOp.CleanOld()
* прибрати харкод зі списку джерел
* вибір джерел для користувача


## tasks

* check .Results after Start/Finish()



## importer

* чекер дірок в імпортованому


## ...

* ПІМ
* тематичний індекс мережі / пошук
* конвертор структур ("семантичний"?)
* панель обробки даних
