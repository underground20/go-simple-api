# Organization management API

Using Go and docker (docker-compose)

### Employees
* Get list employees: `GET /employees`
* Get employee: `GET /employee/:id`
* Create new employee: `POST /employee/add`
* Edit employee: `PUT /employee/:id`
* Delete employee: `DELETE /employee/:id`

Employee info: <code>id, name, sex, age, salary</code>

### Run app
<code>make up</code>