# Organization management API

Using Go and docker (docker-compose)

### Employees
* Get list employees: `GET /employees`
* Get employee: `GET /employee/:id`
* Create new employee: `POST /employee/add`
* Edit employee: `PUT /employee/:id`
* Delete employee: `DELETE /employee/:id`

Employee info: <code>id, name, sex, age, salary, position</code>

Example:
<pre>
{
  "id": 1",
  "name": "Ivan Ivanov",
  "sex": "male",
  "age": 25,
  "salary": 50000,
  "position": "Manager"
}
</pre>

### Departments
* Get department: `GET /department/:id`

Response:
<pre>
{
  "id": 1
  "root_id": 0",
  "name": "Development department",
  "employees": [
    {
      "id": 1
      "name": "Anna Petrova",
      "sex": "female",
      "age": 30,
      "salary": 80000,
      "position": "Senior developer"
    }
  ]
}
</pre>
* Create new department: `POST /department/add`

Body:
<pre>
{
    "id": 2,
    "name": "Development department",
}    
</pre>
* Add employee to existing department: `POST /department/add-employee`

Body
<pre>
{
    "department_id": 2,
    "employee_id": 1,
}
</pre>

Department info: <code>id, root_id, name, employees</code>

### Run app
<code>make up</code>