import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: "5s", target: 1000 },
        { duration: "20s", target: 1000 },
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
            "Authorization": `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwbGF0Zm9ybSI6ImFkbWluIiwic2Vzc2lvbl9pZCI6IjQwMzk1MzYwLTNjOWUtNDY4MC05YTEyLTNjY2Y0NGQzNTM5OSIsInN1YiI6IjFkM2NhMGJmLTgzNWUtNDkwYS1iYTBhLTA0MTQ2ZjVkYTA3ZSIsInVzZXJfcm9sZSI6InN1cGVyYWRtaW4iLCJ1c2VyX3R5cGUiOiJhZG1pbiJ9.-Ocd_nZLGvcOcgEsdXS8NVk9qGBCj0JUJ-HeseKzFXg`,
            "Content-Type": "application/json",
        }
    };


  

    // const resRegister = http.post('http://localhost:8080/v1/user/', Data_Register, registerParams);

    // check(resRegister, {
    //     "status code 201": (r) => r.status === 201
    // });

    const resGetSingleUser = http.get(`http://localhost/v1/user/1d3ca0bf-835e-490a-ba0a-04146f5da07e`, registerParams);

    check(resGetSingleUser, {
        "status code 200": (r) => r.status === 200
    });

    sleep(1)

    // console.log('Register Response:', resRegister.body);
    // console.log('Get Single User Response:', resGetSingleUser.body);
};

