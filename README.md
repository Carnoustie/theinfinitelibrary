Writing social media app for bookworms in go (backend) and typescript (react) frontend. Minimal use of AI assistants and no copy-paste-coding to preserve control of codebase.

The database is of the mysql variety, and assuming mysql is installed, create the database and its schema by running<br>
`mysql -u <your_mysql_username> -p theinfinitelibrary < schemas_theinfinitelibrary.sql`,<br>
after which you will be prompted for your mysql password. The created mysql database is now named *theinfinitelibrary*.

To compile and run the backend, execute `go run .` while standing in the folder `theinfinitelibrary-backend`.

To run the frontend, execute `npm start` while standing in the folder `theinfinitelibrary-frontend`

Frontend runs at an adress hooked to the following domain name:

http://www.theinfinitelibrary.club/ 

# Architecture

## DevOps Architecture

![DevOps workflow](DevOpsArchitecture.png)
