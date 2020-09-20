# Claps.dev

```
 ____     __       ______  ____    ____              ____    ____    __  __    
/\  _`\  /\ \     /\  _  \/\  _`\ /\  _`\           /\  _`\ /\  _`\ /\ \/\ \   
\ \ \/\_\\ \ \    \ \ \L\ \ \ \L\ \ \,\L\_\         \ \ \/\ \ \ \L\_\ \ \ \ \  
 \ \ \/_/_\ \ \  __\ \  __ \ \ ,__/\/_\__ \   _______\ \ \ \ \ \  _\L\ \ \ \ \ 
  \ \ \L\ \\ \ \L\ \\ \ \/\ \ \ \/   /\ \L\ \/\______\\ \ \_\ \ \ \L\ \ \ \_/ \
   \ \____/ \ \____/ \ \_\ \_\ \_\   \ `\____\/______/ \ \____/\ \____/\ `\___/
    \/___/   \/___/   \/_/\/_/\/_/    \/_____/          \/___/  \/___/  `\/__/ 
```

> Help you funding the creators and projects you appreciate with crypto currencies.
## How to run

### Create a Mixin bot

1. Visit https://mixin.one, install and download Mixin Messenger
2. Visit https://developers.mixin.one/dashboard, create a new bot, fill the bot info
3. Copy `client ID`(应用 ID), generate `client secret`(应用密钥) and the `keystore json file`(应用 Session) at "密钥"

### Create a Github OAuth App

1. Visit https://github.com/settings/developers and create a new Github OAuth App
2. Visit https://github.com/settings/tokens and create a new personal token for development

### Config 

1. Configure port number and log information through `conf/application.yml` 
2. Create a file named `.env` like `.env.expmple`
3. Required environment variables:

  ```sh
    GITHUB_CLIENT_ID=YOUR_GITHUB_OAUTH_APP_CLIENT_ID
    GITHUB_CLIENT_SECRET=YOUR_GITHUB_OAUTH_APP_CLIENT_SECRET
    GITHUB_CLIENT_TOKEN=YOUR_GITHUB_OAUTH_APP_CLIENT_TOKEN
    GITHUB_OAUTH_CALLBACK=YOUR_GITHUB_OAUTH_CALLBACK
    
    MIXIN_CLIENT_ID=YOUR_MIXIN_BOT_CLIENT_ID
    MIXIN_CLIENT_CONFIG=PATH_OF_KEYSTORE_FILE
    MIXIN_CLIENT_SECRET=YOUR_MIXIN_BOT_CLIENT_SECRET
    
    DATABASE_ENGINE=YOUR_DATABASE_ENGINE
    DATABASE_HOST=YOUR_DATABASE_HOST
    DATABASE_PORT=YOUR_DATABASE_PORT
    DATABASE_USERNAME=YOUR_DATABASE_USERNAME
    DATABASE_PASSWORD=YOUR_DATABASE_PASSWORD
    DATABASE_DATABASE=YOUR_DATABASE_DATABASE

  ```

### Run

1. go build
