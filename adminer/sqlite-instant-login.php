<?php

class AdminerSQLiteInstantLogin
{
    function credentials()
    {
        $password = get_password();
        return array(SERVER, $_GET["username"], $password);
    }

    function login($login, $password)
    {
        return true;
    }

    function disableInstantLogin()
    {
        $val = $_ENV['ADMINER_DISABLE_INSTANT_LOGIN'] ?? false;
        if (is_string($val) && strtolower($val) === "false") {
            return false;
        }
        return (bool)$val;
    }

    function loginForm()
    {
        if(empty($_GET[DRIVER]) && empty($_GET["username"]) && empty($_GET["db"])): ?>
            <script<?php echo nonce(); ?>>
                var disableInstantLogin = <?php echo json_encode($this->disableInstantLogin()); ?>;
                document.addEventListener("DOMContentLoaded", function(event) {
                    document.querySelector("option[value='sqlite']").selected = true;
                    var s = document.querySelector("input[name='auth[server]']");
                    s.value = "";
                    document.querySelector("input[name='auth[username]']").value = "";
                    document.querySelector("input[name='auth[password]']").value = "";
                    document.querySelector("input[name='auth[db]']").value = "/var/www/html/superdupercloud.db";
                    if (!disableInstantLogin) {
                        s.closest("form").submit()
                    }
                })
            </script><?php
        endif;

        return null;
    }
}

return new AdminerSQLiteInstantLogin();
