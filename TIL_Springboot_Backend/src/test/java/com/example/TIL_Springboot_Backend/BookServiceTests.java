package com.example.TIL_Springboot_Backend;

import jakarta.persistence.EntityNotFoundException;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Assertions;
import org.springframework.boot.test.context.SpringBootTest;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@SpringBootTest
public class BookServiceTests {

    private BookRepository bookRepository = mock(BookRepository.class);
    private UserRepository userRepository = mock(UserRepository.class);

    private BookService bookService = new BookService(bookRepository, userRepository);

    @Test
    void addBookToExistingUser(){
        Long testUserId = 1L;
        String testTitle = "Brave New World";
        String testAuthor = "Aldous Huxley";

        when(userRepository.existsById(testUserId)).thenReturn(true);
        when(bookRepository.save(any())).thenAnswer(i -> i.getArgument(0));

        Book result = bookService.addBookToUser(testTitle, testAuthor, testUserId);
        Assertions.assertEquals(result.getTitle(), testTitle);
        Assertions.assertEquals(result.getAuthor(), testAuthor);
        verify(bookRepository).save(any());
    }

    @Test
    void addBookToNonExistingUser(){
        Long testUserId = 2L;
        String testTitle = "The Mandarins";
        String testAuthor = "Simone de Beauvoir";

        when(userRepository.existsById(2L)).thenReturn(false);

        Assertions.assertThrows(EntityNotFoundException.class, () -> {
            bookService.addBookToUser(testTitle, testAuthor, testUserId);
        });

    }

}
