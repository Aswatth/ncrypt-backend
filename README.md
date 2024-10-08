<h1>Ncrypt</h1>

This is the backend for the password and secret managing application built using <a href="https://gin-gonic.com/">Go</a>
using <a href="https://dgraph.io/docs/badger/" > Badger DB </a> for storing data as key-value pairs.

This acts as a local web-server running on dynamically assigned ports for
the <a href="https://github.com/Aswatth/ncrypt-frontend">desktop application</a>.
JWT authentication has been used to prevent unauthorized access to sensitive data.

<h3>API endpoints:</h3>

<h6>System:</h6>
<table>
  <tr>
    <th>Action</th>
    <th>Path</th>
    <th>Request data</th>
    <th>Description</th>
    <th>Need authentication</th>
  </tr>
  <tr>
    <td>POST</td>
    <td>/system/setup</td>
    <td>
    {
        "master_password": "string",
            "auto_backup_setting": {
              "is_enabled": bool, 
              "backup_location": "string", 
              "backup_file_name": "string" 
            }
    }
    </td>
    <td>Will be used for initial setup. Require master_password and auto_backup_setting data</td>
    <td>No</td>
  </tr>
  <tr>
    <td>POST</td>
    <td>/system/signin</td>
    <td>{"master_password": "string"}</td>
    <td>SIgn into the application</td>
    <td>No</td>
  </tr>
  <tr>
    <td>GET</td>
    <td>/system/generate_password</td>
    <td>-</td>
    <td>Generates random password based on preferences</td>
    <td>No</td>
  </tr>
  <tr>
    <td>POST</td>
    <td>/system/import</td>
    <td>
        {
        "file_name": "string",
        "path": "string",
        "master_password": "string"
        }
    </td>
    <td>Import a given file from the specified path using the master_password to decrypt and load the file</td>
    <td>No</td>
  </tr>
  <tr>
    <td>GET</td>
    <td>system/theme</td>
    <td>-</td>
    <td>Fetches theme(LIGHT, DARK or SYSTEM)</td>
    <td>No</td>
  </tr>
    <tr>
        <td>PUT</td>
        <td>/system/automatic_backup_setting</td>
        <td>
            {
            "auto_backup_setting": {
                          "is_enabled": bool, 
                          "backup_location": "string", 
                          "backup_file_name": "string" 
                        }
            }
        </td>
        <td>Updates auto backup setting information</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>GET</td>
        <td>/system/password_generator_preference</td>
        <td>-</td>
        <td>Fetches password generator preference</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>PUT</td>
        <td>/system/password_generator_preference</td>
        <td>
        {
        "has_digits": bool,
        "has_uppercase": bool,
        "has_special_char": bool,
        "length": int
        }
        </td>
        <td>Update password generator preference</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>GET</td>
        <td>/system/session_duration</td>
        <td>-</td>
        <td>Fetches updates JWT token with extended session</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>PUT</td>
        <td>/system/session_duration</td>
        <td>{"session_duration_in_minutes": int}</td>
        <td>Updates session duration based on given minutes</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>PUT</td>
        <td>/system/theme</td>
        <td>{"theme":"string"}</td>
        <td>Updates theme of the application</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>GET</td>
        <td>/system/data</td>
        <td>-</td>
        <td>Fetches system data</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>POST</td>
        <td>/system/logout</td>
        <td>-</td>
        <td>Logout out of the application</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>POST</td>
        <td>/system/export</td>
        <td>{ "file_name": "string", "path": "string"}</td>
        <td>Export data to specified path with given file name</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>POST</td>
        <td>/system/backup</td>
        <td>-</td>
        <td>Backup data using path and file name sepcified in auto backup setting</td>
        <td>Yes</td>
    </tr>
</table>

<h6>Master password </h6>

<table>
    <tr>
        <th>Action</th>
        <th>Path</th>
        <th>Request data</th>
        <th>Description</th>
        <th>Need authentication</th>
    </tr>
    <tr>
        <td>POST</td>
        <td>/master_password/validate</td>
        <td>{"master_password": "string"}</td>
        <td>Validates the given master password</td>
        <td></td>
    </tr>
    <tr>
        <td>PUT</td>
        <td>/master/password</td>
        <td>{"old_master_password": "string", "new_master_password": "string"}</td>
        <td>Update master password</td>
        <td>Yes</td>
    </tr>
</table>

<h6>Login data</h6>

<table>
    <tr>
        <th>Action</th>
        <th>Path</th>
        <th>Request data</th>
        <th>Description</th>
        <th>Need authentication</th>
    </tr>
    <tr>
        <td>POST</td>
        <td>/login</td>
        <td>{"name": "string", "url": "string", "attributes: {"isFavourite": bool, "requireMasterPassword", bool}, "accounts": [{"username": "string", "password": "string"}]}</td>
        <td>Add new login data with corresponding accounts</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>GET</td>
        <td>/login</td>
        <td>-</td>
        <td>If login/name=? is not specified the call will return all data else will return matching record</td>
    </tr>
    <tr>
        <td>GET</td>
        <td>/login/:name?username=?</td>
        <td>-</td>
        <td>Fetch decrypted account password for given login data and username</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>DELETE</td>
        <td>/login/:name</td>
        <td>-</td>
        <td>Delete login data</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>PUT</td>
        <td>/login/:name</td>
        <td>{"name": "string", "url": "string", "attributes: {"isFavourite": bool, "requireMasterPassword", bool}, "accounts": [{"username": "string", "password": "string"}]}</td>
        <td>Update login data</td>
        <td>Yes</td>
    </tr>
</table>

<h6>Notes:</h6>

<table>
    <tr>
        <th>Action</th>
        <th>Path</th>
        <th>Request data</th>
        <th>Description</th>
        <th>Need authentication</th>
    </tr>
    <tr>
        <td>POST</td>
        <td>/note</td>
        <td>{"created_date_time": "string", "title": "string", "content": "string", "attributes: {"isFavourite": bool, "requireMasterPassword", bool}</td>
        <td>Add new note</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>GET</td>
        <td>/note</td>
        <td>-</td>
        <td>If note/created_date_time=? is not specified the call will return all notes else will return matching note</td>
    </tr>
    <tr>
        <td>GET</td>
        <td>/login/:created_date_time</td>
        <td>-</td>
        <td>Fetch decrypted content for given created_date_time</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>DELETE</td>
        <td>/note/:created_date_time</td>
        <td>-</td>
        <td>Delete note</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>PUT</td>
        <td>/note/:created_date_time</td>
        <td>{"created_date_time": "string", "title": "string", "content": "string", "attributes: {"isFavourite": bool, "requireMasterPassword", bool}</td>
        <td>Update login data</td>
        <td>Yes</td>
    </tr>
</table>

Features:

- Import and export of login data and notes happen in parallel with the help go-routines.
- On master password update, all encrypted data are re-encrypted in parallel.
- Runs on dynamically assigned ports.

To run tests please comment lines 51-65 in system_service.go to prevent UI instances for each test.
