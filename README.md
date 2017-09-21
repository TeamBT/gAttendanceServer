# gAttendanceServer

URL: https://arcane-spire-19269.herokuapp.com

## How to use

Get all students: ```/student```

Get specific student: ```/student/show?id=STUDENT_ID```

Check a student in with id: ```/student/update?id=STUDENT_ID&checkedIn=BOOLEAN&excused=BOOLEAN```
                    
Check a student in with rfid: ```/student/update?rfid=STUDENT_RFID```
                      
Delete a student: ```/student/delete?id=STUDENT_ID```

Reset database so everyone is marked absent: ```/reset```

To create a student make a post request with form data containing name and rfid number to ```/student/create```
