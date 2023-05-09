# Installation
1. Install PostGres on your local machine. I have used the one that is avaiable at https://www.postgresql.org/download/
2. Instead of using SQL shell, I have used PGAdmin (4) to create databases and tables.
3. Once you start the server, it should by default work on port 5432. If you wish to change it, you would need to specify that during installation. I am not certain if it can be modified later. But, there are conf files in the installation directory which should allow for this change to be done later as well.
5. Install Go using Chocolatey (https://www.digitalocean.com/community/tutorials/how-to-install-go-and-set-up-a-local-programming-environment-on-windows-10)
6. Be very careful about where you install the Go environment. The most critical part of it is the $GOPATH. The reason behind this is any further modules that you install would refer to the $GOPATH for installation.
7. Set-up VSCode

APIs:
----

I have used PostMan to test the APIs.
in the code, I have written 2 completely functional APIs whereas the 3rd API for Login is work in progress.
Note: All the Go lang server is running on local host. I have hardcoded port 10000 in the code but it is subject to change depending on the availability of an open port.

The request and responses are as follows:
1. API - /users/
2[. GET: hhtp://localhost:10000/users/](http://localhost:10000/users/)
{"Type":"success","Data":[{"userid":1,"username":"user1","city":"Lanham","email":"user1@user.com"},{"userid":2,"username":"user2","city":"Shuvojit","email":"shuvo@user.com"},{"userid":56,"username":"Shuvojit","city":"Asansol","email":"something@gmail.com"},{"userid":1000,"username":"ShuvojitR","city":"","email":""},{"userid":100011,"username":"ShuvojitRR","city":"","email":"jjja"}],"Message":"Done!"}


2. API - /SignUp/
3. http://localhost:10000/SignUp/?userID=5000&username=Shuv&password=KL&city=Kolkata&email=konami@gmail.com
{"Type":"success","Data":null,"Message":"User has been signed up."}
