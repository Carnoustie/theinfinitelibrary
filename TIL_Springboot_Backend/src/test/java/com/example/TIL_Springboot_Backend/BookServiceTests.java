package com.example.TIL_Springboot_Backend;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Assertions;
import org.springframework.boot.test.context.SpringBootTest;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

@SpringBootTest
public class BookServiceTests {

    private BookRepository bookRepository = mock(BookRepository.class);


    @Test
    void createBookTest(){
        String title = "Brave New World";
        String author = "Aldous Huxley";
        Book b = new Book(title, author);

        Assertions.assertEquals(b.getTitle(), title);
    }

}
