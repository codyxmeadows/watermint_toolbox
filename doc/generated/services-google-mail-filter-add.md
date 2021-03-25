# services google mail filter add

Add a filter. 

# Security

`watermint toolbox` stores credentials into the file system. That is located at below path:

| OS      | Path                                                               |
|---------|--------------------------------------------------------------------|
| Windows | `%HOMEPATH%\.toolbox\secrets` (e.g. C:\Users\bob\.toolbox\secrets) |
| macOS   | `$HOME/.toolbox/secrets` (e.g. /Users/bob/.toolbox/secrets)        |
| Linux   | `$HOME/.toolbox/secrets` (e.g. /home/bob/.toolbox/secrets)         |

Please do not share those files to anyone including Dropbox support.
You can delete those files after use if you want to remove it. If you want to make sure removal of credentials, revoke application access from setting or the admin console.

Please see below help article for more detail:
* Google: https://support.google.com/accounts/answer/3466521

## Auth scopes

| Description                                      |
|--------------------------------------------------|
| Gmail: View and modify but not delete your email |
| Gmail: Manage your basic mail settings           |

# Authorization

For the first run, `tbx` will ask you an authentication with your Google account. Please copy the link and paste it into
your browser. Then proceed to authorization. After authorization, Dropbox will show you an authorization code. Please
copy that code and paste it to the `tbx`.
```

watermint toolbox xx.x.xxx
==========================

© 2016-2020 Takayuki Okazaki
Licensed under open source licenses. Use the `license` command for more detail.

1. Visit the URL for the auth dialogue:

https://accounts.google.com/o/oauth2/auth?client_id=xxxxxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Flocalhost%3A7800%2Fconnect%2Fauth&response_type=code&state=xxxxxxxx

2. Click 'Allow' (you might have to login first):
3. Copy the authorisation code:
Enter the authorisation code
```

# Usage

This document uses the Desktop folder for command example.

## Run

Windows:
```
cd $HOME\Desktop
.\tbx.exe services google mail filter add 
```

macOS, Linux:
```
$HOME/Desktop/tbx services google mail filter add 
```

Note for macOS Catalina 10.15 or above: macOS verifies Developer identity. Currently, `tbx` is not ready for it. Please select "Cancel" on the first dialogue. Then please proceed "System Preference", then open "Security & Privacy", select "General" tab.
You may find the message like:
> "tbx" was blocked from use because it is not from an identified developer.

And you may find the button "Allow Anyway". Please hit the button with your risk. At second run, please hit button "Open" on the dialogue.

## Options:

| Option                      | Description                                                                                                    | Default |
|-----------------------------|----------------------------------------------------------------------------------------------------------------|---------|
| `-add-label-if-not-exist`   | Create a label if it is not exist.                                                                             | false   |
| `-add-labels`               | List of labels to add to the message, separated by ','.                                                        |         |
| `-criteria-exclude-chats`   | Whether the response should exclude chats.                                                                     | false   |
| `-criteria-from`            | The sender's display name or email address.                                                                    |         |
| `-criteria-has-attachment`  | Messages that have any attachment.                                                                             | false   |
| `-criteria-negated-query`   | Only return messages not matching the specified query.                                                         |         |
| `-criteria-no-attachment`   | Messages that does not have any attachment.                                                                    | false   |
| `-criteria-query`           | Only return messages matching the specified query.                                                             |         |
| `-criteria-size`            | The size of the entire RFC822 message in bytes, including all headers and attachments.                         | 0       |
| `-criteria-size-comparison` | How the message size in bytes should be in relation to the size field.                                         |         |
| `-criteria-to`              | The recipient's display name or email address. Includes recipients in the "to", "cc", and "bcc" header fields. |         |
| `-forward`                  | Email address that the message should be forwarded to.                                                         |         |
| `-peer`                     | Account alias                                                                                                  | default |
| `-remove-labels`            | List of labels to remove from the message, separated by ','.                                                   |         |
| `-user-id`                  | The user's email address. The special value me can be used to indicate the authenticated user.                 | me      |

## Common options:

| Option            | Description                                                                               | Default              |
|-------------------|-------------------------------------------------------------------------------------------|----------------------|
| `-auto-open`      | Auto open URL or artifact folder                                                          | false                |
| `-bandwidth-kb`   | Bandwidth limit in K bytes per sec for upload/download content. 0 for unlimited           | 0                    |
| `-budget-memory`  | Memory budget (limits some feature to reduce memory footprint)                            | normal               |
| `-budget-storage` | Storage budget (limits logs or some feature to reduce storage usage)                      | normal               |
| `-concurrency`    | Maximum concurrency for running operation                                                 | Number of processors |
| `-debug`          | Enable debug mode                                                                         | false                |
| `-experiment`     | Enable experimental feature(s).                                                           |                      |
| `-lang`           | Display language                                                                          | auto                 |
| `-output`         | Output format (none/text/markdown/json)                                                   | text                 |
| `-proxy`          | HTTP/HTTPS proxy (hostname:port). Please specify `DIRECT` if you want skip setting proxy. |                      |
| `-quiet`          | Suppress non-error messages, and make output readable by a machine (JSON format)          | false                |
| `-secure`         | Do not store tokens into a file                                                           | false                |
| `-verbose`        | Show current operations for more detail.                                                  | false                |
| `-workspace`      | Workspace path                                                                            |                      |

# Results

Report file path will be displayed last line of the command line output. If you missed command line output, please see path below. [job-id] will be the date/time of the run. Please see the latest job-id.

| OS      | Path pattern                                | Example                                                |
|---------|---------------------------------------------|--------------------------------------------------------|
| Windows | `%HOMEPATH%\.toolbox\jobs\[job-id]\reports` | C:\Users\bob\.toolbox\jobs\20190909-115959.597\reports |
| macOS   | `$HOME/.toolbox/jobs/[job-id]/reports`      | /Users/bob/.toolbox/jobs/20190909-115959.597/reports   |
| Linux   | `$HOME/.toolbox/jobs/[job-id]/reports`      | /home/bob/.toolbox/jobs/20190909-115959.597/reports    |

## Report: filter

Gmail filter
The command will generate a report in three different formats. `filter.csv`, `filter.json`, and `filter.xlsx`.

| Column                 | Description                                                              |
|------------------------|--------------------------------------------------------------------------|
| id                     | Filter Id                                                                |
| criteria_from          | Filter criteria: The sender's display name or email address.             |
| criteria_to            | Filter criteria: The recipient's display name or email address.          |
| criteria_subject       | Filter criteria: Case-insensitive phrase found in the message's subject. |
| criteria_query         | Filter criteria: Only return messages matching the specified query.      |
| criteria_negated_query | Filter criteria: Only return messages not matching the specified query.  |

If you run with `-budget-memory low` option, the command will generate only JSON format report.

In case of a report become large, a report in `.xlsx` format will be split into several chunks like follows; `filter_0000.xlsx`, `filter_0001.xlsx`, `filter_0002.xlsx`, ...

# Proxy configuration

The executable automatically detects your proxy configuration from the environment. However, if you got an error or you want to specify explicitly, please add -proxy option, like -proxy hostname:port. Currently, the executable doesn't support proxies which require authentication.

