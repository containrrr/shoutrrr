# Lark

## URL Format

!!! info ""
    lark://__`host`__/__`token`__?[secret=__`secret`__]
    
--8<-- "docs/services/lark/config.md"

## Create Custom Bot in Lark

Official Documents: [Link](https://open.larksuite.com/document/client-docs/bot-v3/add-custom-bot)

1. Invite custom bot join group.

    a. Enter the target group, click the `More` button in the upper right corner of the group, and then click `Settings`.
    
    b. On the right-side `Settings`, click on `Group Bot`.

    c. Click `Add a Bot` on the `Group Bot`.

    b. In `Add Bot` dialog box, find the `Custom Bot` and add it.

    e. Set the name and description of the custom robot, and click `Add`.

2. Get the webhook address of the custom robot and click `Finish`.

## Get Host and Token of Custom Bot

If you are using `Lark`, then the `Host` is `open.larksuite.com`.

If you are using `Feishu` or `飞书`, then the `Host` is `open.feishu.cn`.

`Token` is the last part of the webhook address. For example, if the webhook address is `https://open.larksuite.com/open-apis/bot/v2/hook/xxxxxxxxxxxxxxxxx`, then the token corresponds to `xxxxxxxxxxxxxxxxx`.

## Get Secret of Custom Bot

1. In the group settings, open the bot list, find the custom bot and click on it to enter the configuration page.

2. In the `Security Settings`, select `Signature Verification`.

3. Click `Copy` to copy the secret.

4. Click `Save` to make the configuration take effect.