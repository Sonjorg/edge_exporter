<?php
    // Prerequisite - Setting up the HTTP Options
        // initialize curl
        $curlHandle = curl_init();
        
        // uncomment for debug porposes
        //curl_setopt($curlHandle, CURLOPT_VERBOSE, true)
        
        // set the appropriate timeout
        curl_setopt($curlHandle, CURLOPT_TIMEOUT, 10);
        
        curl_setopt($curlHandle, CURLOPT_HEADER, false);
        // Set so curl_exec returns the result instead of outputting it.
        curl_setopt($curlHandle, CURLOPT_RETURNTRANSFER, true);
        
        // This example blindly accepts any server certificate, without doing any
        // verification as to which CA signed it, and whether or not that CA is trusted.
        
        // For the sake of simplicity, configure cURL to accept any server(peer) certificate
        curl_setopt($curlHandle, CURLOPT_SSL_VERIFYPEER, false);


        // Following captures the details on how libcurl options can be configured to use secure HTTP.
    
        //curl_setopt($curlHandle , CURLOPT_SSL_VERIFYPEER, true);
        
        // 2, means, check that the common name exists and that it matches the host name of the server
        //curl_setopt($curlHandle , CURLOPT_SSL_VERIFYHOST, 2);
        
        // previously downloaded server cert
        //$certLocation = getcwd() . "/CAcerts/sbc_rest.crt";
        //curl_setopt($curlHandle , CURLOPT_CAINFO, $certLocation);




    // How to Acquire a SBC Session Token
        // define an array where we would cache the Session Token Cookie with value
        $cookieArr = array();
        
        // Sonus SBC 1000/2000 REST login resource URL
        $loginResource = "https://10.233.230.11/rest/login";
        
        // set the login resource url in curl
        curl_setopt($curlHandle, CURLOPT_URL, $loginResource );
        
        // setup a callback handler for reading and processing the response header fields
        // that would come in response to REST login resource POST call
        curl_setopt($curlHandle, CURLOPT_HEADERFUNCTION, array($this, 'responseHeaderCallback'));
        
        // tell cURL that we are doing a HTTP POST
        curl_setopt($curlHandle, CURLOPT_POST, true);
        
        // set the login resource POST params. The user must be of the REST access level and created from the WebUI prior to using the API.
        $loginPropsArr = array('Username'=>'admin', 'Password'=>'admin');
        curl_setopt($curlHandle, CURLOPT_POSTFIELDS, http_build_query($loginPropsArr, '', '&'));
        
        // exec the HTTP/REST request
        $response = curl_exec($curlHandle);
        
        // HTTP response header callback function which processes the header fields
        function responseHeaderCallback($curlHandle, $header) {
            if(!strncmp($header, "Set-Cookie:", 11)) {
                $cookiestr = trim(substr($header, 11, -1));
                $cookie = explode(';', $cookiestr);
                $cookie = explode('=', $cookie[0]);
                $cookiename = trim(array_shift($cookie));
                $cookieArr[$cookiename] = trim(implode('=', $cookie));
            }
            return strlen($header);
        }




        // How to Use the Session Token
            // init cURL handle
            $curlHandle = curl_init();
            
            $cookieHeader = '';
            // previously extracted cookies in $cookieArr (above), is used to add the
            // session token in HTTP request header for subsequent REST call
            foreach ($cookieArr as $key=>$value) {
                $cookieHeader .= "$key=$value; ";
            }
            if (!empty($cookieHeader)) {
                curl_setopt($curlHandle, CURLOPT_COOKIE, $cookieHeader);
            }
            
            // set other relevant HTTP option as shows in above section _Setting up the HTTP Options_
            curl_setopt($curlHandle, CURLOPT_HTTPGET, true);
            
            // Sonus SBC 1000/2000 REST system resource URL
            $systemResource = "https://10.233.230.11/rest/system";
            
            // set the system resource url in curl
            curl_setopt($curlHandle, CURLOPT_URL, $systemResource );
            
            // exec the HTTP/REST request
            $response = curl_exec($curlHandle);




        // How to Close a SBC Session
            /*// init cURL handle
            $curlHandle = curl_init();
            
            $cookieHeader = '';
            // previously extracted cookies in $cookieArr (above), is used to add the
            // session token in HTTP request header for this REST call
            foreach ($cookieArr as $key=>$value) {
                $cookieHeader .= "$key=$value; ";
            }
            if (!empty($cookieHeader)) {
                curl_setopt($curlHandle, CURLOPT_COOKIE, $cookieHeader);
            }
            
            // tell cURL that we are doing a HTTP POST
            curl_setopt($curlHandle, CURLOPT_POST, true);
            
            // set other relevant HTTP option as shows in above section _Setting up the HTTP Options_
            
            // Sonus SBC 1000/2000 REST logout resource URL
            $logoutResource = "https://sbc_host_or_ipaddress/rest/logout";
            
            // set the system resource url in curl
            curl_setopt($curlHandle, CURLOPT_URL, $logoutResource);
            // logout resource does not require any POST parameter.
            
            // exec the HTTP/REST request
            $response = curl_exec($curlHandle);
            
            // any subsequent REST call should result in HTTP Error 401 Unauthorized status*/
?>
