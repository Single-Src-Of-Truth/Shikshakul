package com.modulixinfotech.iams.shikshakul.repository;

import com.modulixinfotech.iams.shikshakul.domain.model.RefreshToken;
import com.modulixinfotech.iams.shikshakul.domain.model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface RefreshTokenRepository extends JpaRepository<RefreshToken, UUID> {
    Optional<RefreshToken> findByTokenId(String tokenId);

    Optional<RefreshToken> findByTokenIdAndRevokedFalse(String tokenId);

    List<RefreshToken> findByUserAndRevokedFalse(User user);

    List<RefreshToken> findByUserAndDeletedFalse(User user);

    @Modifying
    @Query("UPDATE RefreshToken rt SET rt.revoked = true, rt.revokedAt = :revokedAt WHERE rt.user = :user AND rt.revoked = false")
    void revokeAllUserTokens(@Param("user") User user, @Param("revokedAt") LocalDateTime revokedAt);

    @Modifying
    @Query("UPDATE RefreshToken rt SET rt.revoked = true, rt.revokedAt = :revokedAt WHERE rt.tokenId = :tokenId")
    void revokeByTokenId(@Param("tokenId") String tokenId, @Param("revokedAt") LocalDateTime revokedAt);

    List<RefreshToken> findByExpiresAtBeforeAndDeletedFalse(LocalDateTime dateTime);
}