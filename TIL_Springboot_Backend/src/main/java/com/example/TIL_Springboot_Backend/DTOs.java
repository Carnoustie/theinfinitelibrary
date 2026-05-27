package com.example.TIL_Springboot_Backend;

public class DTOs {
    record UserRequest(String username, String password){}
    record BookRequest(String title, String author){};

}
