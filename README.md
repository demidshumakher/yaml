
Приколы:
- не валидный yaml - мы падаем, чтобы жизнь медом не казалась. 
- Даты парсятся только YYYY-MM-DD
- 99.99999% что есть баги, а при багах мы падаем
- многострочные строки очень часто все ломают, они практически не работают
- сложные ключи для словаря, вообще не поддерживаются, потому что сложные типы в большинстве форматах все равно не поддерживается + тут итак все на костылях и святом духе держится 
- отсутствует set
- числа есть только в десятичной CC, просто пишите по-людски и все норм будет)
- отсутствуют какие-то там глобальные теги, уже не помню, что это такое
- может быть много багов, если использовать спецсимволы(:,|< и тд) в строке без кавычек
- тесты вообще для слабаков, главное верить, что оно работает

Что реализованно:
- конвертер из yaml в json
- конвертер из yaml в toml

В теории работает:
- типы данных из json
- sequence 
- key/value pair
- частично поддерживаются теги
    - Обрабатываются только !!str и !!binary теги, остальные еще не добавлены
    - глобальные теги на весь файл не работают
- якоря
- комментарии 
- Можно в быстро добавить поддержку потока файлов, потому что все для этого есть, но оно не обрабатывается, аналогично с концом потока
