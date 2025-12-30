package com.modulixinfotech.iams.shikshakul.repository;

import com.modulixinfotech.iams.shikshakul.domain.model.Role.RoleName;
import com.modulixinfotech.iams.shikshakul.domain.model.Tenant;
import com.modulixinfotech.iams.shikshakul.domain.model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface UserRepository extends JpaRepository<User, UUID> {
    Optional<User> findByEmailAndTenant(String email, Tenant tenant);

    Optional<User> findByMobileAndTenant(String mobile, Tenant tenant);

    Optional<User> findByIdAndDeletedFalse(UUID id);

    Optional<User> findByEmailAndDeletedFalse(String email);

    Optional<User> findByEmailAndTenantAndDeletedFalse(String email, Tenant tenant);

    List<User> findByTenantAndDeletedFalse(Tenant tenant);

    @Query("SELECT u FROM User u JOIN u.roles r WHERE r.name = :roleName AND u.tenant = :tenant AND u.deleted = false")
    List<User> findByRoleAndTenant(@Param("roleName") RoleName roleName, @Param("tenant") Tenant tenant);

    boolean existsByEmailAndTenant(String email, Tenant tenant);

    boolean existsByMobileAndTenant(String mobile, Tenant tenant);
}