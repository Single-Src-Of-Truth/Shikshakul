package com.modulixinfotech.iams.shikshakul.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
@Slf4j
public class ScheduledTasks {

    private final InviteService inviteService;
    private final TokenService tokenService;

    @Scheduled(cron = "0 0 * * * *")
    public void expireOldInvites() {
        log.info("Running scheduled task: Expire old invites");
        try {
            inviteService.expireOldInvites();
        } catch (Exception e) {
            log.error("Error expiring old invites", e);
        }
    }

    @Scheduled(cron = "0 0 2 * * *")
    public void cleanupExpiredTokens() {
        log.info("Running scheduled task: Cleanup expired tokens");
        try {
            tokenService.cleanupExpiredTokens();
        } catch (Exception e) {
            log.error("Error cleaning up expired tokens", e);
        }
    }
}