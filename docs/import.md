# Import

You can import data from various sources.

## YNAB 4

You can import a YNAB 4 budget directly into Envelope Zero. Please read this whole section carefully before you do so.

### Notable differences

The following is **not yet supported** in Envelope Zero and will therefore be ignored:

- Recurring transactions. They will be implemented with [milestone #5](https://github.com/envelope-zero/backend/milestone/5) and import supported with https://github.com/envelope-zero/backend/issues/379.
- Payee rename rules are not yet supported in Envelope Zero (https://github.com/envelope-zero/backend/issues/373)

The following **work differently** on Envelope Zero:

- Date formatting. While YNAB 4 does date formatting per Budget, in Envelope Zero, the formatting is decided by the browser (https://github.com/envelope-zero/frontend/issues/145) or by the users configuration (https://github.com/envelope-zero/backend/issues/33)
- Transactions always need to have a source and destination. Transactions that do not have a Payee set in YNAB 4 will be imported with the opposing account as „YNAB 4 Import - No Payee“.
- Transactions can not have an amount of 0 - if no money was moved, no transaction is needed. Any transaction with an amount of 0 will be ignored during the import.
- In YNAB 4, setting the overspend handling for an Envelope to “Subtract it from next month's category balance” affects all future months for that envelope until you set it back to ”Subtract it from next month's ’Available to Budget‘”. In Envelope Zero, setting overspend handling to `AFFECT_ENVELOPE` (where it is subtracted from the balance of the envelope in the next month) only affects the month it is set on. All months where no overspend handling is configured default to subtract it from the “Available to Budget” sum of the month after. If you set overspend handling to “Subtract it from next month's category balance” for a category in YNAB 4 and never changed it back on any of the following months, during the import, EZ will keep this behaviour for all months until the one you're executing the import in. For following months, you will have to set it manually.

### How to import

YNAB 4 saves all data in a file that can be imported directly to Envelope Zero. The file has the name `Budget.yfull`, you can find it by doing the following:

1. Go to the directory where your YNAB 4 budget file is saved. This file is actually a directory, but it claims to be a file. You can enter this directory with
1. On MacOS: Right click -> Show package content
1. On Windows: ?
1. Open the `Budget.ymeta` file with any text editor. This file tells us which directory to look into next. Check the `relativeDataFolderName` string. Open this directory
1. In the directory, you will find another directory with a UUID as name, for example `F90E864E-8D96-4E0E-A723-776CEEB1C2F0`. Open this directory, too.
1. In there, you will find a Budget.yfull file. This is the one you need, copy it to e.g. your Desktop.
1. Navigate to http://example.com/api/docs/index.html#/Import/post_v1_import_ynab4. Replace `example.com` with the URL of your Envelope Zero instance.
1. Click the "Try it out" button on the top right
1. Select your file and a Budget Name (it must be a name that does not exist yet), then click "Execute".
1. Check the "Server response" section that appeared below. It should have a "Code" of 204. If it does not, check the box at the right for the error message.
1. All done!
