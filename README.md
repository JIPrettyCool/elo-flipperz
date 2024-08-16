# **Backend of Elo Flipperz**

It's a game where to players flip a coin to gain elo

## **Why?**

I wanted to learn Go lang so started this project

# **I Want To Try This!**

### Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`API_KEY`

after that run this project with

```bash
  go run main.go
```

## **Usage**

- **Register**
POST to `http://localhost:8080/register`
```
{
    "username": "BenDover",
    "password": "Handsomeaf"
}
```

- **Login**
GET `http://localhost:8080/login` to get queue token
```
{
    "username": "BenDover",
    "password": "Handsomeaf"
}
```
- **Queue**
POST to `http://localhost:8080/queue` with 2 different users to start match
with the headers
`Authorization: Bearer token_you_got_with_login`

- **Leaderboard**
Visit `http://localhost:8080/leaderboard` to see Elo of players


## License

[GNU AGPLv3](https://choosealicense.com/licenses/agpl-3.0/)


## Authors

- [@jiprettycool](https://www.github.com/JIPrettyCool)