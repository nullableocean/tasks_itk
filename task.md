Описать User API в proto-файле и реализовать холостую логику, то есть можно просто печатать входные данные на экран 
        ◦ Create - ручка создания нового пользователя в системе. 
            ▪ Request 
                • string name
                • string email
                • string password
                • string password_confirm
                • enum role (требуется создать enum из двух вариантов user и admin)
            ▪ Response 
                • int64 id
        ◦ Get - ручка получения информации о пользователе по его идентификатору. 
            ▪ Request 
                • int64 id
            ▪ Response 
                • int64 id
                • string name
                • string email
                • enum role
                • Timestamp created_at
                • Timestamp updated_at
        ◦ Update - ручка обновления информации о пользователе по его идентификатору. 
            ▪ Request 
                • int64 id
                • StringValue name
                • StringValue email
            ▪ Response 
                • empty
        ◦ Delete - удаление пользователя из системы по его идентификатору. 
            ▪ Request 
                • int64 id
            ▪ Response 
                • Empty
Описать Chat API в proto-файле и реализовать холостую логику, то есть можно просто печатать входные данные на экран 
        ◦ Create - ручка создания нового чата. 
            ▪ Request 
                • []string usernames
            ▪ Response 
                • int64 id
        ◦ Delete - удаление чата из системы по его идентификатору. 
            ▪ Request 
                • int64 id
            ▪ Response 
                • empty
        ◦ SendMessage - ручка отправки сообщения на сервер. 
            ▪ Request 
                • string from
                • string text
                • Timestamp timestamp
            ▪ Response 
                • empty