// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
// <reference path="../support/index.d.ts" />

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

/**
 * Note : This test requires the demo plugin tar file under fixtures folder.
 * Download :
 * make dist latest master and copy to ./e2e/cypress/fixtures/com.mattermost.demo-plugin-0.9.0.tar.gz
 */

import {verifyEphemeralMessage} from 'mattermost-webapp/e2e/cypress/integration/integrations/builtin_commands/helper';

describe('Interact with Hello App features', () => {
    const pluginID = Cypress.config('pluginID');
    const pluginFile = Cypress.config('pluginFile');

    before(() => {
        cy.apiAdminLogin();
        cy.visit('/');

        const config = {
            ServiceSettings: {
                EnableOAuthServiceProvider: true,
            }
        };

        cy.apiUpdateConfig(config);

        // cy.apiRemovePluginById(pluginID);

        // cy.apiUploadPlugin(pluginFile);
        cy.apiEnablePluginById(pluginID);

        // # Clean plugin state
        executeCommand('/apps debug-clean');

        // // # Install HTTP Hello App
        executeCommand('/apps debug-install-http');

        // // # Submit the install modal
        cy.get('#interactiveDialogSubmit').click();
    });

    beforeEach(() => {
        // Extra post before each test to verify later posts
        makePost('break');
    });

    after(() => {
        // cy.apiRemovePluginById(pluginID);
    });

    it('Channel header > Create post', () => {
        cy.get('.channel-header button#send').click();
        fillAndSubmitModal();

        verifyEphemeralMessage('Successfully sent survey');
    });

    it('Command > Open Modal', () => {
        executeCommand('/http-hello message-modal --message hey-modal');
        fillAndSubmitModal();

        verifyEphemeralMessage('Successfully sent survey');
    });

    it('Command > Submit', () => {
        executeCommand('/http-hello message --user @sysadmin --message hey-submit');

        verifyEphemeralMessage('Successfully sent survey');
    });

    it('Post dropdown > Submit', () => {
        cy.getLastPostId().then((postId) => {
            cy.uiClickPostDropdownMenu(postId, 'Survey myself');
        });

        verifyEphemeralMessage('Successfully sent survey');
    });

    it('Post dropdown > Modal', () => {
        cy.getLastPostId().then((postId) => {
            cy.uiClickPostDropdownMenu(postId, 'Survey a user');
        });
        fillAndSubmitModal();

        verifyEphemeralMessage('Successfully sent survey');
    });
});

const fillAndSubmitModal = () => {
    cy.get('.modal-body input[placeholder="enter user ID or @user"]').click().type('sysadmin');
    cy.get('#suggestionList img[alt="sysadmin profile image"]').first().click();

    cy.get('.modal-body textarea').click().type('Heres the message');
    cy.get('#appsModalSubmit').click();
};

const executeCommand = (commandStr: string) => {
    cy.get('#post_textbox').clear().type(commandStr);
    cy.wait(200);
    cy.get('#post_textbox').type('{enter}');
};

const makePost = (commandStr: string) => {
    cy.get('#post_textbox').clear().type(commandStr);
    cy.get('#post_textbox').type('{enter}');
};
