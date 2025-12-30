package com.modulixinfotech.iams.shikshakul.service;

import com.modulixinfotech.iams.shikshakul.config.AppProperties;
import com.modulixinfotech.iams.shikshakul.domain.model.RefreshToken;
import com.modulixinfotech.iams.shikshakul.domain.model.User;
import com.modulixinfotech.iams.shikshakul.exception.generic.ResourceNotFoundException;
import com.modulixinfotech.iams.shikshakul.repository.RefreshTokenRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Slf4j
public class TokenService {

    private final RefreshTokenRepository refreshTokenRepository;
    private final AppProperties appProperties;
    private final AuditService auditService;

    @Transactional
    public RefreshToken createRefreshToken(User user, String deviceInfo, String ipAddress, String userAgent) {
        String tokenId = UUID.randomUUID().toString();
        LocalDateTime expiresAt = LocalDateTime.now()
                .plusSeconds(appProperties.getOauth2().getRefreshTokenTtl());

        RefreshToken refreshToken = RefreshToken.builder()
                .user(user)
                .tokenId(tokenId)
                .expiresAt(expiresAt)
                .revoked(false)
                .deviceInfo(deviceInfo)
                .ipAddress(ipAddress)
                .userAgent(userAgent)
                .lastUsedAt(LocalDateTime.now())
                .build();

        RefreshToken saved = refreshTokenRepository.save(refreshToken);
        log.info("Created refresh token for user: {}", user.getEmail());
        return saved;
    }

    @Transactional(readOnly = true)
    public RefreshToken getRefreshTokenByTokenId(String tokenId) {
        return refreshTokenRepository.findByTokenIdAndRevokedFalse(tokenId)
                .orElseThrow(() -> new ResourceNotFoundException("Refresh token not found or revoked"));
    }

    @Transactional
    public RefreshToken rotateRefreshToken(String oldTokenId, String deviceInfo, String ipAddress, String userAgent) {
        RefreshToken oldToken = getRefreshTokenByTokenId(oldTokenId);

        if (oldToken.getExpiresAt().isBefore(LocalDateTime.now())) {
            throw new IllegalArgumentException("Refresh token expired");
        }

        refreshTokenRepository.revokeByTokenId(oldTokenId, LocalDateTime.now());

        RefreshToken newToken = createRefreshToken(oldToken.getUser(), deviceInfo, ipAddress, userAgent);

        auditService.logAction(oldToken.getUser(), oldToken.getUser().getTenant(),
                "REFRESH_TOKEN_ROTATED", "RefreshToken", newToken.getId(),
                "Refresh token rotated", true, null);

        log.info("Rotated refresh token for user: {}", oldToken.getUser().getEmail());
        return newToken;
    }

    @Transactional
    public void updateLastUsed(String tokenId) {
        RefreshToken token = getRefreshTokenByTokenId(tokenId);
        token.setLastUsedAt(LocalDateTime.now());
        refreshTokenRepository.save(token);
    }

    @Transactional
    public void revokeRefreshToken(String tokenId) {
        refreshTokenRepository.revokeByTokenId(tokenId, LocalDateTime.now());
        log.info("Revoked refresh token: {}", tokenId);
    }

    @Transactional
    public void revokeAllUserTokens(User user) {
        refreshTokenRepository.revokeAllUserTokens(user, LocalDateTime.now());
        auditService.logAction(user, user.getTenant(), "ALL_TOKENS_REVOKED",
                "RefreshToken", null, "All refresh tokens revoked", true, null);
        log.info("Revoked all tokens for user: {}", user.getEmail());
    }

    @Transactional(readOnly = true)
    public List<RefreshToken> getActiveTokensForUser(User user) {
        return refreshTokenRepository.findByUserAndRevokedFalse(user);
    }

    @Transactional
    public void cleanupExpiredTokens() {
        List<RefreshToken> expiredTokens = refreshTokenRepository
                .findByExpiresAtBeforeAndDeletedFalse(LocalDateTime.now());

        for (RefreshToken token : expiredTokens) {
            token.setDeleted(true);
            refreshTokenRepository.save(token);
        }

        if (!expiredTokens.isEmpty()) {
            log.info("Cleaned up {} expired refresh tokens", expiredTokens.size());
        }
    }
}