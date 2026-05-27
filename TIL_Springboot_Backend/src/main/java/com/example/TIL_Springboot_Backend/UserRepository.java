package com.example.TIL_Springboot_Backend;

import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface UserRepository extends JpaRepository<User, Long> {
    List<User> getUserByUsername(String username);
}
