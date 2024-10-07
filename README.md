
<h1>Ncrypt</h1>

This is the backend for the password and secret managing application built using <a href="https://gin-gonic.com/">Go</a> using <a href="https://dgraph.io/docs/badger/" > Badger DB </a> for storing data as key-value pairs.

This acts as a local web-server running on dynamically assigned ports for the <a href="https://github.com/Aswatth/ncrypt-frontend">desktop application</a>.
JWT authentication has been been used to prevent unauthorized access to sensitive datas.

API enpoints:

<h3>System</h3>
<table>
  <tr>
    <th>Path</th>
    <th>Request data</th>
    <th>Description</th>
    <th>Need authentication</th>
  </tr>
  <tr>
    <td>/system/setup</td>
    <td>
      "master_password": string, // Master password needed for initial setup
        "auto_backup_setting": {
          "is_enabled": bool, // Toggle auto backup
          "backup_location": string, // Folder location where auto-backup file should be saved
          "backup_file_name": string // File name to use for auto-backup which will appende with date-time stamp. Should not be empty if is_enabled if _true_
        }
    </td>
    <td>Will be used for initial setup. Require master_password and auto_backup_setting data</td>
    <td>No</td>
  </tr>
</table>

