## Getting Started
1. Clone the repository
>   git clone https://github.com/Hemiun/fedbox.git

And switch to ecommerce branch
>git checkout ecommerce

2. Copy file ".env.dist" in the root project directory to ".env" in the same directory

3. Edit the following row in the ".env" file
> FEDBOX_STORAGE=fs
> FEDBOX_HOSTNAME={yourhost:yourport}
> FEDBOX_LISTEN={yourhost:yourport}
> FEDBOX_HTTPS=false

yourhost:yourport = probably will be "localhost:4000"
For docker FEDBOX_LISTEN  must be an external interface (0.0.0.0:4000, no localhost!)

4. Check your HOSTNAME env. It should be the same as FEDBOX_HOSTNAME

5. Build the project
   make download all

as a result a couple of files will be created in the bin directory

6. Creat the storage structure
 >  ./bin/fedboxctl bootstrap

For fs storage you must call this command twice(!) because of a bug in Fedbox.

7. Create oAuth client
>  ./bin/fedboxctl oauth client add -redirectUri http://{yourhost:yourport}/ping

the output will be something like that
> {"level":"debug","path":"/home/user/fedbox/fs/dev","time":"2023-03-30T17:58:12+03:00","message":"Updated Application: http://fedbox.local/actors/82167514-d2c2-4567-a9e0-58d7ea1e8e83"}
> Client ID: 82167514-d2c2-4567-a9e0-58d7ea1e8e83

we need to write down our clientID ("82167514-d2c2-4567-a9e0-58d7ea1e8e83" in the example above)
8. Create a super-user-actor
>   ./bin/fedboxctl pub actor add -attributedTo=http://{yourhost:yourport}/actors/{clientID}

   In the attributedTo param you should specify clientID from previous action "82167514-d2c2-4567-a9e0-58d7ea1e8e83"
   The command asks for username and password
Example:
>./bin/fedboxctl pub actor add -attributedTo=http://localhost:4000/actors/82167514-d2c2-4567-a9e0-58d7ea1e8e83

finally the need write down userID we created
> {"level":"debug","path":"/home/user/fedbox/fedbox/fs/dev","time":"2023-03-30T18:03:31+03:00","message":"Updated Person: http://localhost:4000/actors/7ed5ef56-165b-4089-942d-cd90de510a4c"}
>Added "Person" [admin]: http://localhost:4000/actors/7ed5ef56-165b-4089-942d-cd90de510a4c
>{"level":"debug","path":"/fedbox/fs/dev","time":"2023-04-10T12:29:43Z","message":"Updated Person: http://localhost:4000/actors/7ed5ef56-165b-4089-942d-cd90de510a4c"}

9. Create ordinary user-actor
   In the attributedTo param you should specify ID for super-user-actor from previous action "http://localhost:4000/actors/7ed5ef56-165b-4089-942d-cd90de510a4c"

> ./bin/fedboxctl pub actor add -attributedTo=http://localhost:4000/actors/7ed5ef56-165b-4089-942d-cd90de510a4c

write down user id from the output

10. Run the server
    ./bin/fedbox -env=dev

11. Check if server is running http://{yourhost:yourport}/ping

12. Postman
    Install Postman
    Import "ecommerce" postman collection from "ecommerce.postman_collection.json" project file (press "Import" button)

Than click to "ecommerce" element and go the "Variables" page

On "Variables" page specify the following:
url={yourhost:yourport}
usr=id of super-user-actor from step 7
usr_pass=password of super-user-actor
cln = id of the client from step 6
cln_pass = password of the client
token = leave it empty, it's for internal use

13. oAuth Token
    Token requesting script executes automatically with every request. The result token will be placed in "token" variable.
    To find the script click to "ecommerce" element and go the "Pre request Script" page.

14. Requests
    In "ecommerce" collection we can find the following folders

product - product related requests
fedbox - original Fedbox endpoints
auth - getting oAuth token
user - user related requests

We can run  any request in any order. However, some of them requires ProductID or UserID to be known to work correctly.





