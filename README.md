
# Tx-Notifier bot for Telegram

## Description

This bot is dedicated to tracking the transactions taking place with certain crypto wallets. 

Due to the fact that it was created for internal business purposes, the application logic is rather narrow and not intended for use by the general public (the code does not assume that the application will be used by a large number of users). The lack of a database is due to the absence of a lot of user data, so a small .json file is used as the data store.

## Usage

The application has two main roles: recipients and one possible administrator. Users don't have any rights except to receive messages from the bot. The administrator's ID is predefined in the environment before the deployment. After launching, the administrator needs to configure the bot: define the addresses of tracked wallets, define the recipients of notifications. After entering the configuration, the bot can be launched.

## Commands

`/help` - Show this help message.

`/status` - Get information about recipients, notification status, last transaction block ID, and tracked wallet address.

`/add_user <parameter>` - Add a new user to the notification list.

`/delete_user <parameter>` - Remove a user from the notification list.

`/add_wallet <parameter>` - Add a new wallet to monitor.

`/delete_wallet <parameter>` - Remove a wallet from monitoring.

`/start_notifications` - Start transaction notifications.

`/stop_notifications` - Stop transaction notifications.