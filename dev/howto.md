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
>   make download all  

as a result a couple of files will be created in the bin directory

6. For app initialization (storage creation, oauth client creation, superuser creation) call bootstrap initApp command: 
 > ./bin/fedboxctl bootstrap initApp -secret=<...> -redirectUri=<...>
 
Please, save data that will be print as result  

7. Run the server
    ./bin/fedbox -env=dev

8. Check if server is running http://{yourhost:yourport}/ping

9. Postman
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

10. oAuth Token
    Token requesting script executes automatically with every request. The result token will be placed in "token" variable.
    To find the script click to "ecommerce" element and go the "Pre request Script" page.

11. Requests
    In "ecommerce" collection we can find the following folders

product - product related requests
fedbox - original Fedbox endpoints
auth - getting oAuth token
user - user related requests

We can run  any request in any order. However, some of them requires ProductID or UserID to be known to work correctly.

## Using SMTP

Email notification feature uses an external SMTP server. For testing or development purposes, any personal mail account will be sufficient. Take the following steps to use yandex.mail as an example:

1. Create an account:

    - Go to mail.yandex.ru
    - Create new account in the standard way.
    - Log in
    - Go to "account management" > "security page".
    - Click "App passwords" in the "Access to your data" section on the bottom of the page.
    - Create and store somewhere the Yandex Mail service password.

    Configuration steps are mail provider dependent and probably different for Gmail, etc.

2. Specify the following parameters in the ".env" configuration file:

    > FEDBOX_SMTP_HOST=smtp.yandex.ru
    > FEDBOX_SMTP_PORT=465
    > FEDBOX_SMTP_USER=my_yandex_user_name
    > FEDBOX_SMTP_PASS=my_app_password
    > FEDBOX_SMTP_FROM=ActivityPub eCommerce admin <my_yandex_user_name@yandex.ru>

3. Test:

    - Run Postman and open "dev/fedbox.postman_collection.json".
    - Go to the "mail" folder of the collection.
    - Specify a valid recipient email and name on the params tab.
    - Click send. You should see the "Mail was sent" message in the response.
