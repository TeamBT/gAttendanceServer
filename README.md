# gAttendanceServer

URL: https://arcane-spire-19269.herokuapp.com

## How to use

Get all students: (GET) ```/student```

Get specific student: (GET) ```/student/show?id=STUDENT_ID```

Check a student in with id: (PUT) ```/student/update?id=STUDENT_ID&checkedIn=BOOLEAN&excused=BOOLEAN```

Check a student in with rfid: (PUT) ```/student/update?rfid=STUDENT_RFID```

Delete a student: (DELETE) ```/student/delete?id=STUDENT_ID```

Reset database so everyone is marked absent: (GET) ```/reset```

To create a student: (POST) ```/student/create?name=STUDENT_NAME&rfid=STUDENT_RFID```
