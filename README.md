# Login/Logoff Graph (LLG)

### Log format

- user1.log
    ```
    login;user1;PC (0.0.0.0);01.01.1970;00:00
    logoff;user1;PC (0.0.0.0);01.01.1970;00:01
    ```

- user2.log
    ```
    login;user2;PC (0.0.0.0);01.01.1970;00:00
    logoff;user2;PC (0.0.0.0);01.01.1970;00:01
    ```
  
### Config format
```json
{
  "path": "/path/to/data",
  "selected_users": 3
}
```

### Run command
```bash
./llg.exe
```