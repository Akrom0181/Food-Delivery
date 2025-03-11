import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: "5s", target: 150 },
        { duration: "20s", target: 150 },
        { duration: "5s", target: 0 },
    ],
};

export default () => {
    let Data_Login = JSON.stringify({
        email: "akromjonotaboyev@gmail.com",
        password: "Akrom2005",
        platform: "admin",
    });

    // let uniqueId = Math.random().toString(36).substring(2, 8);
    // let emailDomain = "gmail.com";

    // let Data_Register = JSON.stringify({
    //     full_name: "Test",
    //     user_type: "user",
    //     user_role: "user",
    //     username: "testusername",
    //     email: `test+${uniqueId}@${emailDomain}`,
    //     profile_picture: `${uniqueId}`,
    //     status: "inverify",
    //     password: "1234",
    //     gender: "male",
    // });
    

    // let loginParams = {
    //     headers: {
    //         "Content-Type": "application/json"
    //     }
    // };

    

    let registerParams = {
        headers: {
            "Authorization": `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwbGF0Zm9ybSI6Im1vYmlsZSIsInNlc3Npb25faWQiOiJiNzNiNDczOC1mZDVlLTQwM2UtYmI4Mi0wNDcwOWFkZTY3YmEiLCJzdWIiOiIxZmJkMGUxYy04OGQ3LTRmODAtOTE1ZS01Mjk1ZDY1OGJkOTAiLCJ1c2VyX3JvbGUiOiJ1c2VyIiwidXNlcl90eXBlIjoidXNlciJ9.BcPBzvoFL3Z_0di3Ue1zP13kNk7m-wB-h8A9blKfoN4`,
            "Content-Type": "application/json",
        }
    };


  

    // const resRegister = http.post('http://localhost:8080/v1/user/', Data_Register, registerParams);

    // check(resRegister, {
    //     "status code 201": (r) => r.status === 201
    // });

    const resGetSingleUser = http.get(`http://localhost:9090/v1/report/list?page=1&limit=10`, registerParams);

    check(resGetSingleUser, {
        "status code 200": (r) => r.status === 200
    });

    sleep(1)

    // console.log('Register Response:', resRegister.body);
    // console.log('Get Single User Response:', resGetSingleUser.body);
};

