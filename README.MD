# **Детальное описание архитектуры сервиса**

## **1. Общий обзор**

Данная архитектура является микросервисной и обеспечивает функционал для работы с музыкой: загрузка, сохранение, проигрывание, поиск по названию и метаданным, управление альбомами. Все запросы обрабатываются через **API Gateway**, который распределяет запросы к нужным сервисам.
### Архитектура системы

![System Architecture](images/architecture.png) 

### **Главные функции**

- Загрузка аудиофайлов.
- Создание альбомов и добавление в них песен.
- Проигрывание песен.
- Поиск по названию и векторный поиск по аудиоданным.
- Возможность модерации альбомов или песен.

## **2. Описание компонентов**

### **2.1. API Gateway**
- Выполняет роль единой точки входа в систему.
- Обеспечивает аутентификацию и выдачу JWT-токенов.
- Маршрутизирует запросы к сервисам **Audio Service** и **Albums Service**.

### **2.2. Audio Service**
- Предоставляет API для загрузки, проигрывания и поиска аудио.
- Сохраняет аудиофайлы в хранилище **S3**.
- Записывает метаданные о песнях в базу **PostgreSQL**.
- Отправляет сообщения в **Kafka** для асинхронной обработки аудио сервисом **Audio Processing Service**.

### **2.3. Albums Service**
- Отвечает за управление альбомами и привязкой песен к альбомам.
- Хранит информацию об альбомах в базе **PostgreSQL**.
- Взаимодействует с **Stream Service** для потоковой передачи аудио.

### **2.4. Audio Processing Service**
- Обрабатывает аудиофайлы, поступающие от **Audio Service**.
- Выполняет векторную обработку для поиска по похожим аудиофайлам.
- Взаимодействует с базой **PostgreSQL** и хранилищем **S3**.

### **2.5. Stream Service**
- Отвечает за потоковое воспроизведение песен и альбомов.
- Получает данные от **Albums Service** и **Audio Service**.

### **2.6. PostgreSQL**
- Реляционная база данных, в которой хранятся:
  - Данные о пользователях.
  - Песни и их метаданные.
  - Альбомы.
  - Теги.

### **2.7. S3**
- Хранилище объектов, где сохраняются аудиофайлы и изображения альбомов.

### **2.8. Kafka**
- Используется для асинхронного обмена сообщениями между сервисами.
- Гарантирует, что сообщения о добавленных песнях будут обработаны **Audio Processing Service**.

## **3. Функциональные требования**
- **Прослушивание музыки** – пользователи могут воспроизводить загруженные аудиофайлы.
- **Создание альбомов** – пользователи могут создавать и редактировать альбомы.
- **Добавление песен** – пользователи могут добавлять новые песни в систему.
- **Поиск по названию** – реализован векторный поиск по названиям песен.
- **Моя волна** – функция, предоставляющая персонализированный список песен и альбомов.

## **4. Основные потоки**

### **4.1. Поток добавления песни**
1. Пользователь отправляет аудиофайл через **API Gateway** в **Audio Service**.
2. **Audio Service** сохраняет файл в **S3** и добавляет метаданные о песне в **PostgreSQL**.
3. Далее отправляет сообщение в **Kafka**, которое обрабатывается **Audio Processing Service**.
4. **Audio Processing Service** создает векторное представление аудиофайла и добавляет его в базу данных.

![System Architecture](images/flow_add_song.png) 


### **4.2. Поток создания альбома**
1. Пользователь отправляет запрос на создание альбома через **API Gateway** в **Albums Service**.
2. **Albums Service** сохраняет информацию об альбоме в базе данных и возвращает пользователю идентификатор созданного альбома.

### **4.3. Поток аудио поиска**
1. Пользователь отправляет аудиофайл для поиска.
2. **Audio Service** передает файл в **Audio Processing Service**.
3. **Audio Processing Service** выполняет векторный поиск по базе данных и возвращает найденные песни.

![System Architecture](images/flow_song_find.png) 

## **5. Технологии**
- **API Gateway** – для маршрутизации и аутентификации.
- **PostgreSQL** – реляционная база данных для хранения информации о пользователях, песнях и альбомах.
- **S3** – хранилище объектов для аудио и изображений.
- **Kafka** – для асинхронной коммуникации между сервисами.
- **gRPC** – для взаимодействия между микросервисами.
- **JWT** – для аутентификации пользователей.