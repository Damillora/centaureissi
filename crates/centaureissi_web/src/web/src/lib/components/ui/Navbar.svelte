<script>
    import { token } from "$lib/stores";
    import { isTokenExpired, getUsernameFromToken } from "$lib/login-check";

    let menu_shown = $state(false);

    let loggedIn = $state(false);
    let username = $state("");
    token.subscribe((value) => {
        loggedIn = !isTokenExpired(value);
        username = getUsernameFromToken(value);
    });

    const toggleMenu = () => {
        menu_shown = !menu_shown;
    };
</script>

<nav class="navbar is-primary" aria-label="main navigation">
    <div class="navbar-brand">
        <a class="navbar-item" href="/">
            <!-- <img class="image is-32x32" alt="centaureissi logo" src={logo} /> -->
            <strong>centaureissi</strong>
        </a>

        <a
            href={"#"}
            onclick={toggleMenu}
            role="button"
            class="navbar-burger"
            aria-label="menu"
            aria-expanded="false"
        >
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
        </a>
    </div>

    <div class="navbar-menu" class:is-active={menu_shown}>
        <div class="navbar-start">
            <a href="/search" class="navbar-item">Search</a>
        </div>

        <div class="navbar-end">
            <div class="navbar-item has-dropdown is-hoverable">
                <a href={"#"} class="navbar-link">
                    {loggedIn ? username : "Guest"}
                </a>
                {#if loggedIn}
                    <div class="navbar-dropdown is-right">
                        <a href="/stats" class="navbar-item">
                            Stats
                        </a>
                        <a href="/user/profile" class="navbar-item">
                            Profile
                        </a>
                        <a href="/user/password" class="navbar-item">
                            Change Password
                        </a>
                        <hr class="navbar-divider" />
                        <a href="/auth/logout" class="navbar-item">Log out</a>
                    </div>
                {:else}
                    <div class="navbar-dropdown">
                        <a href="/auth/register" class="navbar-item">
                            Register
                        </a>
                        <a href="/auth/login" class="navbar-item">Log in</a>
                    </div>
                {/if}
            </div>
        </div>
    </div>
</nav>
