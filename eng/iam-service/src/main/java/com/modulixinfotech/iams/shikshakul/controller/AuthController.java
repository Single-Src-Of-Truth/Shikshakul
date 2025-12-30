package com.modulixinfotech.iams.shikshakul.controller;

import com.modulixinfotech.iams.shikshakul.domain.model.Invite;
import com.modulixinfotech.iams.shikshakul.mapper.InviteDTO;
import com.modulixinfotech.iams.shikshakul.mapper.common.ApiResponse;
import com.modulixinfotech.iams.shikshakul.service.AuthService;
import com.modulixinfotech.iams.shikshakul.service.InviteService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

import static com.modulixinfotech.iams.shikshakul.mapper.request.AuthRequest.RegisterRequest;
import static com.modulixinfotech.iams.shikshakul.mapper.response.AuthResponse.RegisterResponse;

@RestController
@RequestMapping("/api/auth")
@RequiredArgsConstructor
@Slf4j
public class AuthController {

    private final AuthService authService;
    private final InviteService inviteService;

    @PostMapping("/register")
    public ResponseEntity<ApiResponse<RegisterResponse>> register(
            @Valid @RequestBody RegisterRequest request) {
        log.info("Registration request received for email: {}", request.getEmail());

        RegisterResponse response = authService.registerUser(request);

        return ResponseEntity
                .status(HttpStatus.CREATED)
                .body(ApiResponse.success("Registration successful", response));
    }

    @GetMapping("/invite/{token}")
    public ResponseEntity<ApiResponse<InviteDTO>> validateInvite(@PathVariable UUID token) {
        log.info("Validating invite token: {}", token);

        Invite invite = inviteService.validateInviteToken(token);
        InviteDTO inviteDTO = InviteDTO.fromEntity(invite);

        return ResponseEntity.ok(ApiResponse.success("Valid invite", inviteDTO));
    }

    @GetMapping("/check-email")
    public ResponseEntity<ApiResponse<Boolean>> checkEmailAvailability(
            @RequestParam String email,
            @RequestParam String tenantDomain) {
        log.info("Checking email availability: {} for tenant: {}", email, tenantDomain);

        return ResponseEntity.ok(ApiResponse.success(false));
    }
}