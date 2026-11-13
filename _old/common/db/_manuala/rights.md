# Перевірка/встановлення прав на записах в БД

## 0. Права на перегляд запису. 

Це робиться тільки по полю rView. rOwner ні до чого — система не повинна дозволяти виставляти таке rView, на яке власник запису не має прав. Отже, для контролю перегляду — розпаковувати mamagers не треба.



## 1. Перевірка при створенні запису і формування прав на наступні операції з записом.

func SetRights(
   identity          Identity, 
   ctrl              Controller,
   interfaceManagers Managers, 
   newRecordManagers Managers) (IdentityString, IdentityString, []byte, error)

Створення запису дозволяється тільки якщо користувач має права відповідні до interfaceManagers[rights.Create], инакше — помилка (зокрема, коли interfaceManagers не задано або там нема значення для rights.Create).

Якщо newRecordManagers == nil або там не задано rights.RView і/або rights.ROwner — ці права встановлюються згідно Identity

Всі права, які задано в newRecordManagers повинні бути у самого користувача. Не можна встановити право, якого сам не маєш. Така спроба — помилка.

Якщо всі перевірки пройшли, то функція повертає оновлені rView та rOwner і упаковане поле managers.

Якщо ctrl не задано, права перевіряються тільки на співпадіння IdentityString, якщо ж заданий — то також і на BelongsTo().



## 2. Перевірка/оновлення прав при спробі зміни запису.

Змінювати можна тільки сам запис, а можна змінювати і якісь права доступу до нього. Отже, нехай все це перевіряється
одною функцією:

func CheckAndUpdateRights(
   identity                Identity, 
   ctrl                    Controller, 
   rView, rOwner, managers string, 
   updatedRecordManagers   Managers) (IdentityString, IdentityString, []byte, error)

І правила наступні.

rView, rOwner, managers — це поля запису, "як є". managers треба розпакувати, а rView та rOwner конвертнути в IdentityString. 

updatedRecordManagers (може містити ключі rights.RView, rights.ROwner) треба порівняти з наявними правами —
якщо updatedRecordManagers == nil або в updatedRecordManagers нема ніяких змін щодо поточних прав — "права на зміну прав"
не перевіряємо.

Якщо ж в updatedRecordManagers задана зміна прав, то:
* дозволяємо її тільки якщо identity має права rights.ROwner, инакше — помилка;
* всі права, які задано в updatedRecordManagers повинні бути у самого користувача — не можна встановити право,
якого сам не маєш, така спроба — помилка.
* якщо всі перевірки пройшли успішно, то функція повертає оновлені rView та rOwner і перепаковане поле managers
(в цьому випадку вважаємо, що заодно, користувач має право міняти і все инше в записі, більше перевірок не треба).

Якщо ж в updatedRecordManagers не була задана зміна прав, то:
* слід перевірити права на зміну запису: rights.ROwner або rights.RChange, якщо ця перевірка не успішна — помилка;
* якщо ж перевірка пройшла успішно — функція повертає вихідні значення rView, rOwner, managers без змін.

Якщо ctrl не задано, права перевіряються тільки на співпадіння IdentityString, якщо ж заданий — то також і на BelongsTo().

## 3. Перевірка/оновлення прав при спробі вилучення запису.

func CheckDeleteRights(
   identity         Identity, 
   ctrl             Controller, 
   rOwner, managers string) error

rOwner, managers — це поля запису, "як є". Перевіряємо, чи identity має права rOwner або managers[rights.Delete].  
