package com.modulixinfotech.iams.shikshakul.security.userdetails;

import com.modulixinfotech.iams.shikshakul.domain.model.User;
import com.modulixinfotech.iams.shikshakul.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@RequiredArgsConstructor
@Slf4j
public class CustomUserDetailsService implements UserDetailsService {

    private final UserRepository userRepository;

    @Override
    @Transactional(readOnly = true)
    public UserDetails loadUserByUsername(String email) throws UsernameNotFoundException {
        User user = userRepository.findByEmailAndDeletedFalse(email)
                .orElseThrow(() -> new UsernameNotFoundException("User not found with email: " + email));

        log.debug("Loading user: {}", email);

        return new CustomUserDetails(user);
    }

    @Transactional(readOnly = true)
    public UserDetails loadUserByEmailAndTenantDomain(String email, String tenantDomain)
            throws UsernameNotFoundException {
        User user = userRepository.findByEmailAndDeletedFalse(email)
                .filter(u -> u.getTenant().getDomain().equals(tenantDomain))
                .orElseThrow(() -> new UsernameNotFoundException(
                        "User not found with email: " + email + " in tenant: " + tenantDomain));

        log.debug("Loading user: {} for tenant: {}", email, tenantDomain);

        return new CustomUserDetails(user);
    }
}