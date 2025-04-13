package by.kapinskiy.Discussion.controllers;


import by.kapinskiy.Discussion.DTOs.Requests.NoteRequestDTO;
import by.kapinskiy.Discussion.DTOs.Responses.NoteResponseDTO;
import by.kapinskiy.Discussion.services.NotesService;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/notes")
public class NotesController {
    private final NotesService notesService;

    @Autowired
    public NotesController(NotesService notesService) {
        this.notesService = notesService;
    }


    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public NoteResponseDTO createNote(@RequestBody @Valid NoteRequestDTO noteRequestDTO) {
        return notesService.save(noteRequestDTO);
    }

    @GetMapping
    @ResponseStatus(HttpStatus.OK)
    public List<NoteResponseDTO> getAllNotes() {
        return notesService.findAll();
    }

    @GetMapping("/{id}")
    @ResponseStatus(HttpStatus.OK)
    public NoteResponseDTO getNoteById(@PathVariable Long id) {
        return notesService.findById(id);
    }

    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteNote(@PathVariable long id) {
        notesService.deleteById(id);
    }

    // Non REST version for tests compliance
    @PutMapping
    @ResponseStatus(HttpStatus.OK)
    public NoteResponseDTO updateNote(@RequestBody @Valid NoteRequestDTO noteRequestDTO) {
        return notesService.update(noteRequestDTO);
    }

/*    @GetMapping
    @ResponseStatus(HttpStatus.OK)
    public List<NoteResponseDTO> getAllNotes() {
        return Collections.emptyList();
    }

    @GetMapping("/{id}")
    @ResponseStatus(HttpStatus.OK)
    public NoteResponseDTO getNoteById(@PathVariable Long id) {
        return new NoteResponseDTO();
    }*/
}
